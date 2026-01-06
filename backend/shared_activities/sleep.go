package sharedactivities

import (
	"context"
	"time"
)

func (r *RobotActivities) Sleep(ctx context.Context, durationSeconds int) error {
	// Simulate sleep activity
	time.Sleep(time.Duration(durationSeconds) * time.Second)
	return nil
}
