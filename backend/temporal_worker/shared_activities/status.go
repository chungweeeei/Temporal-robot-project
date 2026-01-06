package sharedactivities

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"go.temporal.io/sdk/activity"
)

func (r *RobotActivities) FetchStatus(ctx context.Context, url string) (string, error) {

	logger := activity.GetLogger(ctx)

	// Get or create connection
	conn, err := r.getConnection(ctx, url)
	if err != nil {
		logger.Error("Failed to get websocket connection", "url", url, "error", err)
		return "", err
	}

	// execute business logic
	data := map[string]int{
		"api_id": 1009,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	commandBytes, err := r.generateCommand(string(dataBytes))
	if err != nil {
		return "", err
	}

	err = conn.WriteMessage(websocket.TextMessage, commandBytes)
	if err != nil {
		logger.Error("Failed to write message, removing connection", "url", url, "error", err)
		r.removeConnection(url)
		return "", err
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		logger.Error("Failed to read message, removing connection", "url", url, "error", err)
		r.removeConnection(url)
		return "", err
	}

	return string(message), nil
}
