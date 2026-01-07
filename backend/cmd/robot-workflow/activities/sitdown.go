package activities

import (
	"context"
	"encoding/json"

	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) Sitdown(ctx context.Context, url string) (string, error) {

	logger := activity.GetLogger(ctx)

	// execute business logic
	data := map[string]int{
		"api_id": RobotMotionControl,
		"action": SitDownActionID,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := ra.Client.CallService(ctx, url, string(dataBytes))
	if err != nil {
		logger.Error("Failed to send sitdown command", "error", err)
		return "", err
	}
	return response, nil
}
