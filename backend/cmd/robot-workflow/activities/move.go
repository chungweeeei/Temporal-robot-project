package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) Move(ctx context.Context, params map[string]interface{}) (string, error) {

	logger := activity.GetLogger(ctx)

	// validataion
	targetX, okX := params["x"].(float64)
	targetY, okY := params["y"].(float64)
	if !okX || !okY {
		return "", fmt.Errorf("invalid parameters for Move activity")
	}

	// Step 1: send move command
	_, err := executeWithHeartbeat(ctx, func() (string, error) {
		data := map[string]interface{}{
			"api_id":      RobotMoveCommandID,
			"mission_id":  uuid.New().String(),
			"x":           targetX,
			"y":           targetY,
			"orientation": 0.0,
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		response, err := ra.Client.CallService(ctx, string(dataBytes))
		if err != nil {
			logger.Error("Failed to send move command", "erorr", err)
			return "", err
		}
		logger.Info("Move command accepted by robot", "response", response)
		return response, nil
	})
	if err != nil {
		return "", err
	}

	// Step 2: Polling until robot move to target
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// TODO: send stop command to robot?
			return "", ctx.Err()
		case <-ticker.C:
			// 1. fetch current robot status
			status, err := ra.GetStatus(ctx)
			if err != nil {
				logger.Warn("Failed to poll robot status", "error", err)
				continue
			}

			if status.Status.Code != 0 {
				return "", fmt.Errorf("robot reported error during move: code %d", status.Status.Code)
			}

			// 2. send heartbeart, tell temporal server we are still alive
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Robot currently at (%f, %f)", status.X, status.Y))
			if !status.IsMoving {
				if targetX == 5.0 && targetY == 5.0 {
					return "", fmt.Errorf("robot move failed")
				}

				return "move completed", nil
			}
		}
	}

}
