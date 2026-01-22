package activities

import (
	"context"
	"encoding/json"

	config "github.com/chungweeeei/Temporal-robot-project/internal/config/activity"
	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) Standup(ctx context.Context, params map[string]interface{}) (string, error) {

	logger := activity.GetLogger(ctx)

	return executeWithHeartbeat(ctx, func() (string, error) {
		data := map[string]int{
			"api_id": config.RobotMotionControlID,
			"action": config.StandUpActionID,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		logger.Info("Call Standup Service")
		return ra.Client.CallService(ctx, "Standup", string(dataBytes))
	})
}
