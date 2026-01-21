package activities

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

type RobotClient struct {
	robotURL string
	Dialer   WSDialer
}

func NewRobotClient(robotIP string) *RobotClient {

	return &RobotClient{
		robotURL: fmt.Sprintf("ws://%s:9090", robotIP),
		Dialer: &DefaultDialer{
			Dialer: websocket.DefaultDialer,
		},
	}
}

func (r *RobotClient) CallService(ctx context.Context, actionType string, data any) (string, error) {

	conn, _, err := r.Dialer.DialContext(ctx, r.robotURL, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to get connection: %v", err)
	}
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-doneCh:
		}
	}()
	defer conn.Close()

	payload, err := generatePayload(actionType, data)
	if err != nil {
		return "", err
	}

	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	type readResult struct {
		data []byte
		err  error
	}
	resultCh := make(chan readResult, 1)

	go func() {
		_, msg, err := conn.ReadMessage()
		resultCh <- readResult{data: msg, err: err}
	}()

	select {
	case res := <-resultCh:
		if res.err != nil {
			return "", fmt.Errorf("read failed: %v", res.err)
		}
		return parseResponse(res.data)
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
