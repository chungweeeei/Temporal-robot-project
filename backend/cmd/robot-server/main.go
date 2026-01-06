package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var request robot.CallServiceRequest
		if err := json.Unmarshal(message, &request); err != nil {
			continue
		}

		response := bot.HandleRequest(request)
		conn.WriteJSON(response)
	}
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Println("Mock Robot Server 啟動於 :9090/")
	http.ListenAndServe(":9090", nil)
}
