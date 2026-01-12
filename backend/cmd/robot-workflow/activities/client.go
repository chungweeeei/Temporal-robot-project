package activities

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

const RobotServiceURL = "ws://localhost:9090"

type RobotClient struct {
	Dialer WSDialer
}

func NewRobotClient() *RobotClient {
	return &RobotClient{
		Dialer: &DefaultDialer{
			Dialer: websocket.DefaultDialer,
		},
	}
}

func (r *RobotClient) CallService(ctx context.Context, data string) (string, error) {
	conn, _, err := r.Dialer.DialContext(ctx, RobotServiceURL, nil)
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

	payload, err := generatePayload(data)
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
		fmt.Println("Context done, cancelling read")
		return "", ctx.Err()
	}
}
