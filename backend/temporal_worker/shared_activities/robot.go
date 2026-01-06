package sharedactivities

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type RobotActivities struct {
	Dialer WSDialer
	// mutex proteccts the conns map
	mu    sync.Mutex
	conns map[string]*CachedConnection
}

// CachedConnection wraps a WSConnection with a mutex for thread-safe writes
type CachedConnection struct {
	conn WSConnection
	mu   sync.Mutex
}

func (c *CachedConnection) WriteMessage(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *CachedConnection) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *CachedConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

func NewRobotActivities() *RobotActivities {
	return &RobotActivities{
		Dialer: &DefaultDialer{
			Dialer: websocket.DefaultDialer,
		},
		conns: make(map[string]*CachedConnection),
	}
}

// GetConnection returns an existing connection or creates a new one
func (r *RobotActivities) getConnection(ctx context.Context, url string) (*CachedConnection, error) {
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

func (r *RobotActivities) removeConnection(url string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if conn, ok := r.conns[url]; ok {
		conn.Close()
		delete(r.conns, url)
	}
}

func (r *RobotActivities) generateCommand(data string) ([]byte, error) {
	msg := map[string]any{
		"op":      "call_service",
		"service": "/api/system",
		"type":    "custom_msgs/srv/Api",
		"args": map[string]string{
			"data": data,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}
