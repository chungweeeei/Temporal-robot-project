package client

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/internal/robot/simulator"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gorilla/websocket"
)

type RobotHandler struct {
	bot *simulator.MockRobot
}

func NewRobotHandler(bot *simulator.MockRobot) *RobotHandler {
	return &RobotHandler{bot: bot}
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
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

func (h *RobotHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	safeConn := &SafeConn{conn: conn}
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
		_, message, err := conn.ReadMessage()
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
			var request pkg.RobotServiceRequest
			if err := json.Unmarshal(message, &request); err != nil {
				continue
			}

			go func(req pkg.RobotServiceRequest) {
				response := h.bot.HandleRequest(req)

				if err := safeConn.WriteJSON(response); err != nil {
					log.Println("Write error:", err)
				}
			}(request)
		case "subscribe":
			var request pkg.RobotTopicRequest
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
				go h.RobotStatusBroadcaster(safeConn, done)
			default:
				log.Println("Unknown topic:", request.Topic)
				close(done)
				delete(activeSubscriptions, request.Topic)
			}
		}
	}
}

func (h *RobotHandler) RobotStatusBroadcaster(conn *SafeConn, done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			status := h.bot.GetRobotStatus()
			if err := conn.WriteJSON(status); err != nil {
				log.Println("Broadcast error", err)
				return
			}
		}
	}
}
