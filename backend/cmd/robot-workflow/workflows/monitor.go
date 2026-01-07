package workflows

import (
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"go.temporal.io/sdk/workflow"
)

const WorkflowId = "robot_monitor_workflow"
const TaskQueueName = "ROBOT_MONITOR_TASK_QUEUE"

func RobotMonitorWorkflow(ctx workflow.Context) error {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	robotURL := "ws://localhost:9090"
	var ra *activities.RobotActivities

	logger := workflow.GetLogger(ctx)

	for {
		var result string
		// 1. 執行 Activity 查詢狀態
		err := workflow.ExecuteActivity(ctx, ra.GetStatus, robotURL).Get(ctx, &result)
		if err != nil {
			logger.Error("Health check failed in this cycle", "error", err)
		} else {
			// 2. 檢查狀態邏輯 (例如：電量過低發警報)
			logger.Info("Receive result", result)
			// if err := json.Unmarshal([]byte(result), &state); err == nil {
			// 	if state.BatteryLevel < 20 {
			// 		logger.Warn("Low battery detected!", "level", state.BatteryLevel)
			// 		// 這裡可以觸發其他的 Notification Activity...
			// 	}
			// }
		}

		// 3. 休眠一段時間 (例如每 1 分鐘檢查一次)
		// 使用 workflow.Sleep 非常高效，不會佔用 Worker 資源
		workflow.Sleep(ctx, 30*time.Second)
	}

}
