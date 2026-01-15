package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
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
		response, err := ra.Client.CallService(ctx, "Move", string(dataBytes))
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
			color.Red("Move activity cancelled, stopping robot monitoring")

			// 1. 建立一個不依賴原本 ctx 的新 context (避免因為 cancelled 而發送失敗)
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 2. 準備 Stop Command 的 Payload
			stopData := map[string]interface{}{
				"api_id":     RobotStopActionID, // 使用 configs.go 中定義的 ID (5000)
				"mission_id": uuid.New().String(),
			}
			stopBytes, _ := json.Marshal(stopData)

			// 3. 發送 Stop 指令
			// 注意：這裡使用 stopCtx 而不是原本的 ctx
			_, err := ra.Client.CallService(stopCtx, "Stop", string(stopBytes))
			if err != nil {
				logger.Error("Failed to send stop command during cancellation", "error", err)
			} else {
				logger.Info("Stop command sent successfully")
			}

			return "", ctx.Err()
		case <-ticker.C:
			status, err := ra.GetStatus(ctx)
			if err != nil {
				logger.Error("Failed to get robot status during move", "error", err)
				continue
			}

			distance := ((status.Pose.Position.X - targetX) * (status.Pose.Position.X - targetX)) + ((status.Pose.Position.Y - targetY) * (status.Pose.Position.Y - targetY))
			if distance < 0.01 {
				return "Robot reached target location", nil
			}
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Robot currently at (%f, %f)", status.Pose.Position.X, status.Pose.Position.Y))
		}
	}

}
