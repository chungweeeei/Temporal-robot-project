package activities

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type CheckTaskResp struct {
	HasUpdate bool   `json:"has_update"`
	TaskID    string `json:"task_id"`
}

func CheckTaskUpdates(ctx context.Context) (*CheckTaskResp, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Checking for task updates via REST API...")

	time.Sleep(500 * time.Millisecond)

	return &CheckTaskResp{
		HasUpdate: true,
		TaskID:    "task_12345",
	}, nil
}
