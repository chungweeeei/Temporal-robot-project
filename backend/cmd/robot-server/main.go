package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-server/robot"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	bot = robot.New()
)

type SafeConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *SafeConn) WriteJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	rawConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer rawConn.Close()

	safeConn := &SafeConn{conn: rawConn}

	for {
		_, message, err := rawConn.ReadMessage()
		if err != nil {
			break
		}

		var request robot.CallServiceRequest
		if err := json.Unmarshal(message, &request); err != nil {
			continue
		}

		go func(req robot.CallServiceRequest) {
			response := bot.HandleRequest(req)

			safeConn.WriteJSON(response)
		}(request)
	}
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Println("Mock Robot Server started on :9090")
	http.ListenAndServe(":9090", nil)
}
