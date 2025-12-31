package workflows

import (
	"time"

	sharedactivities "github.com/chungweeeei/robot/shared_activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RobotActionType string

const (
	StandUp RobotActionType = "standup"
	SitDown RobotActionType = "sitdown"
)

func RobotAction(ctx workflow.Context, action RobotActionType) (string, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    3 * time.Second,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	robotURL := "ws://10.8.140.130:9090"
	var ra *sharedactivities.RobotActivities
	var result string
	var err error

	switch action {
	case StandUp:
		err = workflow.ExecuteActivity(ctx, ra.Standup, robotURL).Get(ctx, &result)
	case SitDown:
		err = workflow.ExecuteActivity(ctx, ra.Sitdown, robotURL).Get(ctx, &result)
	default:
		return "unknown behavior", nil
	}

	if err != nil {
		return "", err
	}

	return result, nil
}
