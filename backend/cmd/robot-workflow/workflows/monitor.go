package workflows

import (
	"fmt"
	"os"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"go.temporal.io/sdk/workflow"
)

func RobotMonitorWorkflow(ctx workflow.Context) error {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	robotIP := os.Getenv("ROBOT_IP")
	if robotIP == "" {
		robotIP = "10.8.140.130"
	}
	fmt.Println("robotIp", robotIP)
	robotURL := fmt.Sprintf("ws://%s:9090", robotIP)
	var ra *activities.RobotActivities

	logger := workflow.GetLogger(ctx)

	for {
		var status activities.RobotStatus
		// 1. Execute Activity to get robot status
		err := workflow.ExecuteActivity(ctx, ra.GetStatus, robotURL).Get(ctx, &status)
		if err != nil {
			logger.Error("Health check failed in this cycle", "error", err)
		} else {
			// 2. Check status logic (e.g., alert if battery is low)
			if status.BatteryLevel < 50 {
				logger.Warn("Low battery detected!", "level", status.BatteryLevel)
				// This is where you can trigger other Notification Activities...
				targetWorkflowID := "robot_action_workflow_001"
				signalName := "low-battery-signal"

				signalFuture := workflow.SignalExternalWorkflow(ctx, targetWorkflowID, "", signalName, status.BatteryLevel)
				if err := signalFuture.Get(ctx, nil); err != nil {
					logger.Error("Failed to send signal to external workflow", "targetID", targetWorkflowID, "error", err)
				} else {
					logger.Info("Signal sent successfully", "targetID", targetWorkflowID)
				}
			}
		}

		// 3. Sleep for a while (e.g., check every 1 minute)
		// Using workflow.Sleep is very efficient and does not occupy Worker resources
		workflow.Sleep(ctx, 30*time.Second)
	}

}
