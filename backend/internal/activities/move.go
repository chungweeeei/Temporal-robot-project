package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) sendStopCommand() {

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopData := map[string]interface{}{
		"api_id": config.RobotStopActionID,
	}
	stopBytes, _ := json.Marshal(stopData)

	_, err := ra.Client.CallService(stopCtx, "Stop", string(stopBytes))
	if err != nil {
		color.Red("Failed to send stop command", "error", err)
	}
}

func (ra *RobotActivities) Move(ctx context.Context, params map[string]interface{}) (string, error) {

	logger := activity.GetLogger(ctx)

	// validataion
	targetX, okX := params["x"].(float64)
	targetY, okY := params["y"].(float64)
	targetOrientation, okO := params["orientation"].(float64)
	if !okX || !okY || !okO {
		return "", fmt.Errorf("invalid parameters for Move activity")
	}

	// Step 1: send move command
	newMissionID := uuid.New().String()
	_, err := executeWithHeartbeat(ctx, func() (string, error) {
		data := map[string]interface{}{
			"api_id":      config.RobotMoveCommandID,
			"mission_id":  newMissionID,
			"x":           targetX,
			"y":           targetY,
			"orientation": targetOrientation * (math.Pi / 180.0), // degree to radian
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		response, err := ra.Client.CallService(ctx, "Move", string(dataBytes))
		if err != nil {
			logger.Error("Failed to send move command", "error", err)
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
			logger.Info("Move activity cancelled, stopping robot.")
			ra.sendStopCommand()
			return "", ctx.Err()

		case <-ticker.C:
			status, err := ra.GetStatus(ctx)

			if err != nil {
				// instantly check context error
				if ctx.Err() != nil {
					logger.Info("Move activity cancelled (from GetStatus error).")
					ra.sendStopCommand()
					return "", ctx.Err()
				}
				logger.Error("Failed to get robot status during move", "error", err)
				continue
			}

			if status.MissionID != newMissionID {
				logger.Info("Waiting for robot to start the move mission", "expected_mission_id", newMissionID, "current_mission_id", status.MissionID)
				continue
			}

			if status.Mission.Code == config.MissionSuccess {
				return fmt.Sprintf("Robot has reached the target location (%.2f, %.2f)", status.Pose.Position.X, status.Pose.Position.Y), nil
			}
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Robot currently at (%f, %f)", status.Pose.Position.X, status.Pose.Position.Y))
		}
	}

}
