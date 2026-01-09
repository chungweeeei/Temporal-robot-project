package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

type MoveRequest struct {
	X           int `json:"x"`
	Y           int `json:"y"`
	Orientation int `json:"orientation"`
}

func (ra *RobotActivities) Move(ctx context.Context, url string, req MoveRequest) (string, error) {

	logger := activity.GetLogger(ctx)

	// Step 1: send move command
	data := map[string]interface{}{
		"api_id":      RobotMoveCommandID,
		"mission_id":  uuid.New().String(),
		"x":           req.X,
		"y":           req.Y,
		"orientation": req.Orientation,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := ra.Client.CallService(ctx, url, string(dataBytes))
	if err != nil {
		logger.Error("Failed to send move command", "error", err)
		return "", err
	}
	logger.Info("Move command accepted by robot", "response", resp)

	// Step 2: Polling until robot move to target
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Move activity has been cancelled")
			// TODO: send stop command to robot?
			return "", ctx.Err()
		case <-ticker.C:
			// 1. fetch current robot status
			status, err := ra.GetStatus(ctx, url)
			if err != nil {
				logger.Warn("Failed to poll robot status", "error", err)
				continue
			}

			if status.Status.Code != 0 {
				return "", fmt.Errorf("robot reported error during move: code %d", status.Status.Code)
			}

			logger.Info("Robot is moving...", "current_x", status.X, "current_y", status.Y)

			// 2. send heartbeart, tell temporal server we are still alive
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Robot currently at (%f, %f)", status.X, status.Y))
			if !status.IsMoving {
				return "move completed", nil
			}

		}
	}

}
