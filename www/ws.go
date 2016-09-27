package www

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type wsError struct {
	Error  string `json:"error"`
	Status uint   `json:"status"`
}

type wsOuterMessage struct {
	Project      string            `json:"project"`
	Branch       string            `json:"branch"`
	Action       string            `json:"action"`
	InnerMessage map[string]string `json:"innerMessage,omitempty"`
}

// WSAction is a type of function that handles a WS message
type WSAction func(conn *websocket.Conn, msg *wsOuterMessage) error

// WS is the struct resprenting a websocket server handler
type WS struct {
	Upgrader *websocket.Upgrader
	Logger   *log.Logger
	Actions  map[string]WSAction
}

// NewWS returns a pointer to a WS instance
func NewWS(bufferSize int, actions map[string]WSAction) *WS {
	return &WS{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  bufferSize,
			WriteBufferSize: bufferSize,
		},
		Logger:  log.New(os.Stdout, "WS_SERVER : ", log.LstdFlags),
		Actions: actions,
	}
}

// HandleUpgrades handles a HTTP upgrade to a websocket
func (ws *WS) HandleUpgrades(res http.ResponseWriter, req *http.Request) {
	conn, connErr := ws.Upgrader.Upgrade(res, req, nil)
	if connErr != nil {
		ws.HandleServerError(res, req, connErr)
	}
	ws.Route(conn)
}

// HandleServerError handles WS server errors by logging them and sending a
// response back of the error itself
func (ws *WS) HandleServerError(res http.ResponseWriter, req *http.Request, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	ws.Logger.Println("ERROR : ", err.Error())
	resp := wsError{
		Error:  err.Error(),
		Status: (uint)(http.StatusInternalServerError),
	}
	encErr := json.NewEncoder(res).Encode(resp)
	if encErr != nil {
		ws.Logger.Println("ERROR : ", encErr.Error())
	}
}

// HandleWSError handles errors once the connection has been upgraded
func (ws *WS) HandleWSError(conn *websocket.Conn, err error) {
	ws.Logger.Println("ERROR : ", err.Error())
	data := wsError{
		Error:  err.Error(),
		Status: (uint)(http.StatusInternalServerError),
	}
	writeErr := conn.WriteJSON(data)
	if writeErr != nil {
		ws.Logger.Println("ERROR : ", writeErr.Error())
	}
}

// Route routes messages between a websocket connection
func (ws *WS) Route(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, msg, msgErr := conn.ReadMessage()
		if msgErr != nil {
			ws.HandleWSError(conn, msgErr)
		}
		recvMsg := &wsOuterMessage{}
		unmarshErr := json.Unmarshal(msg, recvMsg)
		if unmarshErr != nil {
			ws.HandleWSError(conn, unmarshErr)
		}
		ws.Logger.Println("MESSAGE : ", recvMsg)
		if ws.Actions[recvMsg.Action] != nil {
			actErr := ws.Actions[recvMsg.Action](conn, recvMsg)
			if actErr != nil {
				ws.HandleWSError(conn, actErr)
			}
		} else {
			ws.HandleWSError(conn, errors.New(fmt.Sprintf("action %s does not exist", recvMsg.Action)))
		}
	}
}
