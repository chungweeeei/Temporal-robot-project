package activities

import (
	"context"
	"fmt"

	transport "github.com/chungweeeei/Temporal-robot-project/internal/transport/websocket"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gorilla/websocket"
)

type RobotClient struct {
	RobotURL string
	Dialer   transport.WSDialer
}

func NewRobotClient(robotIP string) *RobotClient {
	return &RobotClient{
		RobotURL: fmt.Sprintf("ws://%s:9090", robotIP),
		Dialer: &transport.DefaultDialer{
			Dialer: websocket.DefaultDialer,
		},
	}
}

func (r *RobotClient) CallService(ctx context.Context, actionType pkg.ActivityType, data interface{}) (string, error) {

	// Register Websocket Dialer Connection
	conn, _, err := r.Dialer.DialContext(ctx, r.RobotURL, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to get connection: %v", err)
	}
	defer conn.Close()

	doneCh := make(chan struct{})
	defer close(doneCh)

	payload, err := generatePayload(actionType, data)
	if err != nil {
		return "", fmt.Errorf("Failed to generate payload: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return "", fmt.Errorf("Failed to write message via websocket: %v", err)
	}

	type readResult struct {
		data []byte
		err  error
	}
	resultCh := make(chan readResult, 1)

	// Register goroutine read websocket message
	go func() {
		_, msg, err := conn.ReadMessage()
		resultCh <- readResult{data: msg, err: err}
	}()

	// Wait for either read result or context done
	select {
	case res := <-resultCh:
		if res.err != nil {
			return "", fmt.Errorf("Failed to read message via websocket: %v", res.err)
		}
		return parseResponse(res.data)
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
