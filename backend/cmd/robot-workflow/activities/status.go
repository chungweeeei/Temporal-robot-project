package activities

import (
	"context"
	"encoding/json"

	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) GetStatus(ctx context.Context, url string) (string, error) {

	logger := activity.GetLogger(ctx)

	data := map[string]int{
		"api_id": RobotStatus,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := ra.Client.CallService(ctx, url, string(dataBytes))
	if err != nil {
		logger.Error("Failed to get status", "error", err)
		return "", err
	}

	return response, nil
}
