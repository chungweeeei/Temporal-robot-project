package transport

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
WSConnection is an interface representing a WebSocket connection.
It defines methods for writing and reading messages and closing the connection.
*/
type WSConnection interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
	// SetWriteDeadline(t time.Time) error
	// SetReadDeadline(t time.Time) error
}

/*
WSDialer is an interface representing a WebSocket dialer.
It defines a method for dialing a WebSocket connection with a given context, URL, and request headers.
*/
type WSDialer interface {
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (WSConnection, *http.Response, error)
}

type DefaultDialer struct {
	Dialer *websocket.Dialer
}

func (d *DefaultDialer) DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (WSConnection, *http.Response, error) {
	conn, resp, err := d.Dialer.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, resp, err
	}
	return conn, resp, nil
}
