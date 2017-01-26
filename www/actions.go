package www

import (
	"github.com/gorilla/websocket"

	"github.com/cpg1111/maestro-www/data"
)

func stop(interrupt chan bool) {
	interrupt <- false
}

func pipeWS(conn *websocket.Conn, interrupt chan bool, dataChan chan []byte, errChan chan error) error {
	defer stop(interrupt)
	for {
		select {
		case data := <-dataChan:
			if len(data) == 0 {
				break
			}
			_, err := conn.Write(data)
			if err != nil {
				return err
			}
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewWSActions(dClient *data.Client) map[string]WSAction {
	return map[string]WSAction{
		"watch.project": func(conn *websocket.Conn, msg *wsOuterMessage) error {
			interrupt := make(chan bool)
			dataChan, errChan := dClient.WatchProject(msg.Project, interrupt)
			return pipeWS(conn, interrupt, dataChan, errChan)
		},
		"watch.build": func(conn *websocket.Conn, msg *wsOuterMessage) error {
			interrupt := make(chan bool)
			dataChan, errChan := dClient.WatchOne(msg.Project, msg.Branch, interrupt)
			return pipeWS(conn, interrupt, dataChan, errChan)
		},
	}
}
