package activities

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type RobotClient struct {
	Dialer WSDialer
	// mutex protects the connections map
	mu    sync.Mutex
	conns map[string]*CachedConnection
}

func NewRobotClient() *RobotClient {
	return &RobotClient{
		Dialer: &DefaultDialer{
			Dialer: websocket.DefaultDialer,
		},
		conns: make(map[string]*CachedConnection),
	}
}

func (r *RobotClient) CallService(ctx context.Context, url string, payload string) (string, error) {
	conn, err := r.GetConnection(ctx, url)
	if err != nil {
		return "", fmt.Errorf("Failed to get connection: %v", err)
	}

	reqBytes, err := generateCommand(payload)
	if err != nil {
		return "", err
	}

	respBytes, err := conn.WriteAndRead(ctx, reqBytes)
	if err != nil {
		r.RemoveConnection(url)
		return "", fmt.Errorf("communication failed: %v", err)
	}

	respData, err := parseResponse(respBytes)
	if err != nil {
		return "", err
	}

	return respData, nil
}

func (r *RobotClient) GetConnection(ctx context.Context, url string) (*CachedConnection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if conn, ok := r.conns[url]; ok {
		return conn, nil
	}

	// Dial new connection
	c, _, err := r.Dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	cachedConn := &CachedConnection{conn: c}
	r.conns[url] = cachedConn
	return cachedConn, nil
}

func (r *RobotClient) RemoveConnection(url string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if conn, ok := r.conns[url]; ok {
		conn.Close()
		delete(r.conns, url)
	}
}

type CachedConnection struct {
	mu   sync.Mutex
	conn WSConnection
}

func (c *CachedConnection) WriteAndRead(ctx context.Context, data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// deadline := time.Now().Add(10 * time.Second)
	// 假設 WSConnection 介面有 SetDeadline，如果只有 SetWriteDeadline/SetReadDeadline 需分開設
	// c.conn.SetWriteDeadline(deadline)
	// c.conn.SetReadDeadline(deadline)

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, fmt.Errorf("write failed: %v", err)
	}

	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return msg, nil
}

func (c *CachedConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}
