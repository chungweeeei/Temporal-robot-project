package sharedactivities

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSConnection is an interface representing a WebSocket connection.
// It defines methods for writing messages and closing the connection.
// This abstraction allows for easier testing and flexibility in handling WebSocket connections.
type WSConnection interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// WSDialer is an interface representing a WebSocket dialer.
// It defines a method for dialing a WebSocket connection with a given context, URL, and request headers.
// This abstraction allows for easier testing and flexibility in establishing WebSocket connections.
type WSDialer interface {
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (WSConnection, *http.Response, error)
}

type DefaultDialer struct {
	Dialer *websocket.Dialer
}

// DialContext dials a WebSocket connection using the underlying websocket.Dialer.
func (d *DefaultDialer) DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (WSConnection, *http.Response, error) {
	conn, resp, err := d.Dialer.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, resp, err
	}
	return conn, resp, nil
}
