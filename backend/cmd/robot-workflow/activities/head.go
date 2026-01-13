package activities

import (
	"context"
	"encoding/json"
	"fmt"
)

func (ra *RobotActivities) Head(ctx context.Context, params map[string]interface{}) (string, error) {

	// logger := activity.GetLogger(ctx)
	// logger.Info("Head activity called", "params", params)
	_, ok := params["angle"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid parameters for head activity")
	}

	return executeWithHeartbeat(ctx, func() (string, error) {

		data := map[string]float64{
			"data": 45.0,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return ra.Client.CallService(ctx, "head", string(dataBytes))
	})
}
