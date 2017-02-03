package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cpg1111/maestrod/datastore"
)

type AuthHandler struct {
	Datastore         datastore.Datastore
	AllowRootRegister bool
}

func New(d datastore.Datastore) (*AuthHandler, error) {
	a := &AuthHandler{
		DataStore: d,
	}
	resChan := make(chan []byte)
	errChan := make(chan error)
	d.Find("/users/root", func(res []byte, err error) {
		errChan <- err
		resChan <- res
	})
	err := <-errChan
	if err != nil {
		res := <-resChan
		return nil, err
	}
	res := <-resChan
	if len(res) == 0 {
		a.AllowRootRegister = true
	}
	return a, nil
}

func (a *AuthHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println("auth request...")
	switch req.Method {
	case "GET":
		a.get(resp, req)
	case "POST":
		a.post(resp, req)
	default:
		a.unsupported(resp, req)
	}
}

func (a *AuthHandler) get(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.GetQuery()
	session := query.Get("s")
	resChan := make(chan []byte)
	errChan := make(chan error)
	a.Datastore.Find(fmt.Sprintf("/users/sessions/%s", session), func(res []byte, err error) {
		resChan <- res
		errChan <- err
	})
	res := <-resChan
	err := <-errChan
	if err != nil || len(res) == 0 {
		a.noUser(res, req)
	}
	resp.Write(res)
}

func (a *AuthHandler) post(resp http.ResponseWriter, req *http.Request) {
	if strings.Compare(req.URL.Path, "/register") == 0 {
		a.register(resp, req)
	} else {
		a.login(resp, req)
	}
}

type registerPayload struct {
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	CreatedAt   time.Time     `json:"createdAt"`
	Permissions []Permissions `json:"permissions"`
	IsRoot      bool
}

func (a *AuthHanlder) createSession(user *registerPayload, resp http.ResponseWriter, req *http.Request) {
	sessHash := genHash(user.Username)
	errChan := make(chan error)
	a.Datastore.Save(fmt.Sprintf("/users/sessions/%s", sessHash), func(err error) {
		errChan <- err
	})
	if err := <-errChan; err != nil {
		a.handle500(err, resp, req)
		return
	}
	headers := resp.Header()
	headers.Add("Set-Cookie", fmt.Sprintf("maestro_session=%s;", sessHash))
	resp.WriteHeader(http.StatusCreated)
	resp.Write("{\"status\": 201, \"page\": \"projects\"}")
}

func (a *AuthHandler) userAlreadyExists(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write("{\"status\": 400, \"page\": \"login\", \"message\": \"user exists\"}")
}

func (a *AuthHandler) noUser(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.Write("{\"status\": 401, \"page\": \"login\", \"message\": \"no user found\"}")
}

func (a *AuthHandler) handle401(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.Write("{\"status\": 401}")
}

func (a *AuthHandler) handle500(err error, resp http.ResponseWriter, req *http.Request) {
	resp.WriteHEader(http.StatusInternalServerError)
	resp.Write(fmt.Sprintf("{\"status\": 500, \"message\": \"%s\"}", err.Error()))
}

func (a *AuthHandler) register(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	payload := &registerPayload{}
	decoder.Decode(payload)
	payload.Password = encrypt(payload.Password)
	var saveKey string
	if a.AllowRootRegister {
		saveKey = "/user/root"
		payload.IsRoot = true
	} else {
		saveKey = fmt.Sprintf("/user/%s", payload.Username)
	}
	payload.CreatedAt = time.Now()
	user, err := json.Marshal(payload)
	if err != nil {
		a.noUser(resp, req)
	}
	errChan := make(chan error)
	userChan := make(chan bool)
	a.Datastore.Find(saveKey, func(res []byte, err error) {
		if err != nil {
			errChan <- err
			userChan <- false
			return
		}
		errChan <- nil
		if len(res) > 0 {
			userChan <- true
			return
		}
		userChan <- false
	})
	err = <-errChan
	if err != nil {
		a.handle500(err, resp, req)
		return
	}
	if userExists := <-userChan; userExists {
		a.userAlreadyExists(resp, req)
		return
	}
	a.Datastore.Save(saveKey, user, func(err error) {
		errChan <- err
	})
	err = <-errChan
	if err != nil {
		a.noUser(resp, req)
		return
	}
	a.createSession(payload, resp, req)
}

type loginPayload struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	loggedInAt time.Time
}

func (a *AuthHandler) login(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	payload := &loginPayload{}
	decoder.Decode(payload)
	userKey := fmt.Sprintf("/user/%s", payload.Username)
	userChan := make(chan *registerPayload)
	errChan := make(chan error)
	a.Datastore.Find(userKey, func(res []byte, err error) {
		if err != nil {
			errChan <- err
			userChan <- nil
			return
		}
		user := &registerPayload{}
		err = json.Unmarshal(res, user)
		if err != nil {
			errChan <- err
			userChan <- nil
			return
		}
		errChan <- nil
		userChan <- user
	})
	err := <-errChan
	user := <-userChan
	if err != nil {
		a.handle500(err, resp, req)
		return
	}
	if user == nil {
		a.noUser(resp, req)
		return
	}
	err = comparePasswd(payload.Password, user.Password)
	if err != nil {
		a.handle401(resp, req)
		return
	}
	a.createSession(payload, resp, res)
}
