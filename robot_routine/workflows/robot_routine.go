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
	Sleep   RobotActionType = "sleep"
)

type RobotStep struct {
	Action   RobotActionType `json:"action"`
	Duration time.Duration   `json:"duration,omitempty"`
}

func RobotRoutine(ctx workflow.Context, steps []RobotStep) (string, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	robotURL := "ws://10.8.140.130:9090"
	var ra *sharedactivities.RobotActivities

	for _, step := range steps {
		var result string
		var err error
		switch step.Action {
		case StandUp:
			err = workflow.ExecuteActivity(ctx, ra.Standup, robotURL).Get(ctx, &result)
		case SitDown:
			err = workflow.ExecuteActivity(ctx, ra.Sitdown, robotURL).Get(ctx, &result)
		case Sleep:
			workflow.Sleep(ctx, step.Duration)
		default:
			return "unknown behavior", nil
		}
		if err != nil {
			return "", err
		}
	}

	return "routine completed", nil
}
