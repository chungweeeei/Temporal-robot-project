package workflows

import (
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"go.temporal.io/sdk/workflow"
)

type RobotActionType string

const (
	StandUp RobotActionType = "standup"
	SitDown RobotActionType = "sitdown"
)

func RobotActionWorkflow(ctx workflow.Context, action RobotActionType) (string, error) {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	robotURL := "ws://localhost:9090"
	var ra *activities.RobotActivities
	var result string
	var err error

	switch action {
	case StandUp:
		err = workflow.ExecuteActivity(ctx, ra.Standup, robotURL).Get(ctx, &result)
	default:
		return "unknown behavior", nil
	}

	if err != nil {
		return "", err
	}

	return result, nil
}
