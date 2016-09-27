package www

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

var (
	actions = map[string]WSAction{
		"test": func(conn *websocket.Conn, msg *wsOuterMessage) error {
			return nil
		},
	}
	ws = NewWS(4096, actions)
)

func TestHandleUpgrades(t *testing.T) {
	res := httptest.NewRecorder()
	req, reqErr := http.NewRequest("GET", "ws://127.0.0.1:8080/stream", nil)
	if reqErr != nil {
		t.Error(reqErr)
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "websocket")
	go ws.HandleUpgrades(res, req)
	resp := res.Result()
	for resp == nil {
		resp = res.Result()
	}
	if resp.StatusCode >= 500 {
		var err []byte
		resp.Body.Read(err)
		t.Errorf((string)(err))
	}
}
