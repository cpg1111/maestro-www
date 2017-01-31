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
	DataStore         datastore.Datastore
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

func (a *AuthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("auth request...")
	switch req.Method {
	case "GET":
		a.get(res, req)
	case "POST":
		a.post(res, req)
	default:
		a.unsupported(res, req)
	}
}

func (a *AuthHandler) get(res http.ResponseWriter, req *http.Request) {
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

func (a *AuthHandler) post(res http.ResponseWriter, req *http.Request) {
	if strings.Compare(req.URL.Path, "/register") == 0 {
		a.register(res, req)
	} else {
		a.login(res, req)
	}
}

type RegisterPayload struct {
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	ConfPass    string        `json:"confpass"`
	CreatedAt   time.Time     `json:"createdAt"`
	Permissions []Permissions `json:"permissions"`
	IsRoot      bool
}

func (a *AuthHandler) register(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

}
