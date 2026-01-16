package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-server/robot"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
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

type MessageHeader struct {
	Op    string `json:"op"`
	Topic string `json:"topic"`
	Type  string `json:"type"`
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

	// use map to track active subscriptions and their Cancel Channels
	// Key: Topic Nmae, Value: Cancel Channel for that subscription
	activeSubscriptions := make(map[string]chan struct{})
	var subsMu sync.Mutex

	// Clean up function
	defer func() {
		subsMu.Lock()
		for _, cancel := range activeSubscriptions {
			close(cancel)
		}
		subsMu.Unlock()
	}()

	for {
		_, message, err := rawConn.ReadMessage()
		if err != nil {
			break
		}

		// Step One: Check message header
		var header MessageHeader
		if err := json.Unmarshal(message, &header); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		// Step two: Handle based on Op type
		switch header.Op {
		case "call_service":
			var request pkg.ServiceRequest
			if err := json.Unmarshal(message, &request); err != nil {
				continue
			}

			// 每個 Request 獨立處理，避免阻塞讀取迴圈
			go func(req pkg.ServiceRequest) {
				response := bot.HandleRequest(req)

				if err := safeConn.WriteJSON(response); err != nil {
					log.Println("Write error:", err)
				}
			}(request)
		case "subscribe":
			var request pkg.TopicRequest
			if err := json.Unmarshal(message, &request); err != nil {
				continue
			}

			subsMu.Lock()
			defer subsMu.Unlock()

			if _, exists := activeSubscriptions[request.Topic]; exists {
				log.Println("Already subscribed to:", request.Topic)
				continue
			}

			// trigger broadcaster based on topic
			done := make(chan struct{})
			activeSubscriptions[request.Topic] = done

			switch request.Topic {
			case "/api/info":
				go RobotStatusBroadcaster(safeConn, done)
			default:
				log.Println("Unknown topic:", request.Topic)
				close(done)
				delete(activeSubscriptions, request.Topic)
			}
		}
	}
}

func RobotStatusBroadcaster(conn *SafeConn, done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			status := bot.GetRobotStatus()
			if err := conn.WriteJSON(status); err != nil {
				log.Println("Broadcase error", err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/", wsHandler)
	log.Println("Mock Robot Server started on :9090")
	http.ListenAndServe("localhost:9090", nil)
}
