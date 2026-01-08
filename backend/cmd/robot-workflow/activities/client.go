package activities

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

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

func (r *RobotClient) CallService(ctx context.Context, url string, payload string) (string, error) {
	conn, _, err := r.Dialer.DialContext(ctx, url, nil)
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

	reqBytes, err := generateCommand(payload)
	if err != nil {
		return "", err
	}

	if err := conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
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
