package workflows

import (
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"go.temporal.io/sdk/workflow"
)

func RobotScheduleWorkflow(ctx workflow.Context) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("Schedule workflow started", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Check for Task Updates
	var checkResp activities.CheckTaskResp
	err := workflow.ExecuteActivity(ctx, activities.CheckTaskUpdates).Get(ctx, &checkResp)
	if err != nil {
		logger.Error("Failed to check task updates", "Error", err)
		return err
	}

	if checkResp.HasUpdate {
		logger.Info("Task table updated", "TaskID", checkResp.TaskID)
	} else {
		logger.Info("No task table updates found")
	}

	return nil
}
