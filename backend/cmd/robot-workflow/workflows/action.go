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

	// 1. 建立可取消的 Context (關鍵步驟)
	// 這個 childCtx 會被傳給 Activity，而 cancel 函數可以在背景協程中呼叫/
	childCtx, cancel := workflow.WithCancel(ctx)
	defer cancel()

	// 2. Trigger Background Goroutine 監控外部取消事件
	workflow.Go(ctx, func(gCtx workflow.Context) {
		signalChan := workflow.GetSignalChannel(gCtx, "low-battery-signal")

		// 阻塞等待 Signal
		// 假如不需要讀取具體內容，可以傳 nil
		signalChan.Receive(gCtx, nil)

		// 收到 Signal 後，取消主流程的 Context
		workflow.GetLogger(gCtx).Info("Low battery signal received, cancelling activity...")
		cancel() // 呼叫 cancel 取消 childCtx，進而取消 Activity
	})

	robotURL := "ws://localhost:9090"
	var ra *activities.RobotActivities
	var result string
	var err error

	// 3. 執行 Activity (注意這裡使用的是 childCtx)
	// 如果上面的協程執行了 cancel()，這裡會立即收到 CanceledError
	switch action {
	case StandUp:
		err = workflow.ExecuteActivity(childCtx, ra.Standup, robotURL).Get(childCtx, &result)
	default:
		return "unknown behavior", nil
	}

	if err != nil {
		return "", err
	}

	return result, nil
}
