package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) Head(ctx context.Context, params map[string]interface{}) (string, error) {

	logger := activity.GetLogger(ctx)

	angle, ok := params["angle"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid parameters for head activity")
	}

	return executeWithHeartbeat(ctx, func() (string, error) {

		data := map[string]float64{
			"data": angle,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		logger.Info("Call Head Service", "angle", angle)

		return ra.Client.CallService(ctx, string(dataBytes))
	})
}
