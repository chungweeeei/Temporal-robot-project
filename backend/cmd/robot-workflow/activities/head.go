package activities

import (
	"context"
	"fmt"
)

func (ra *RobotActivities) Head(ctx context.Context, params map[string]interface{}) (string, error) {

	angle, ok := params["angle"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid parameters for head activity")
	}

	return executeWithHeartbeat(ctx, func() (string, error) {
		return ra.Client.CallService(ctx, "Head", angle)
	})
}
