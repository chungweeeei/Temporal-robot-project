package activities

import (
	"context"
	"encoding/json"
	"time"
)

func (ra *RobotActivities) Standup(ctx context.Context, url string) (string, error) {

	data := map[string]int{
		"api_id": RobotMotionControl,
		"action": StandUpActionID,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := ra.Client.CallService(ctx, url, string(dataBytes))
	if err != nil {
		return "", err
	}

	time.Sleep(5 * time.Minute)

	return response, nil
}
