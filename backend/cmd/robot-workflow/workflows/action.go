package workflows

import (
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RobotActionType string

const (
	StandUp RobotActionType = "standup"
	SitDown RobotActionType = "sitdown"
	Move    RobotActionType = "move"
)

type RobotWorkflowRequest struct {
	Action     RobotActionType         `json:"action"`
	MoveTarget *activities.MoveRequest `json:"move_target,omitempty"`
}

func RobotActionWorkflow(ctx workflow.Context, req RobotWorkflowRequest) (string, error) {

	logger := workflow.GetLogger(ctx)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// 1. 建立可取消的 Context (關鍵步驟)
	// 這個 childCtx 會被傳給 Activity，而 cancel 函數可以在背景協程中呼叫/
	childCtx, cancel := workflow.WithCancel(ctx)
	defer cancel()

	// 2. Trigger Background Goroutine 監控外部取消事件
	workflow.Go(ctx, func(gCtx workflow.Context) {
		signalChan := workflow.GetSignalChannel(gCtx, "low-battery-signal")
		logger.Info("Waiting for low-battery-signal...")
		// 阻塞等待 Signal
		// 假如不需要讀取具體內容，可以傳 nil
		signalChan.Receive(gCtx, nil)
		logger.Info("signal received! Calling cancel()...")

		// 收到 Signal 後，取消主流程的 Context
		cancel() // 呼叫 cancel 取消 childCtx，進而取消 Activity
		logger.Info("Cancel called.")
	})

	robotURL := "ws://localhost:9090"
	var ra *activities.RobotActivities
	var result string
	var err error

	// 3. 執行 Activity (注意這裡使用的是 childCtx)
	// 如果上面的協程執行了 cancel()，這裡會立即收到 CanceledError
	switch req.Action {
	case StandUp:
		future := workflow.ExecuteActivity(childCtx, ra.Standup, robotURL)
		err = future.Get(ctx, &result)
	case SitDown:
		future := workflow.ExecuteActivity(childCtx, ra.Sitdown, robotURL)
		err = future.Get(ctx, &result)
	case Move:
		future := workflow.ExecuteActivity(childCtx, ra.Move, robotURL, *req.MoveTarget)
		err = future.Get(ctx, &result)
	default:
		return "unknown behavior", nil
	}

	if err != nil {

		if temporal.IsCanceledError(err) {
			logger.Info("Workflow canceled due to low battery, starting recovery...")

			// 關鍵步驟：創建 Disconnected Context
			// 這樣即使父 Context被取消，這個 Activity 仍能執行
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)

			// 執行補償 Activity: SitDown
			// 這裡你可以設定適合 SitDown 的 ActivityOptions (如果需要的話)
			// 例如 SitDown 比較緊急，timeout 可以設短一點
			sitdownOptions := workflow.ActivityOptions{
				StartToCloseTimeout: 5 * time.Minute,
				HeartbeatTimeout:    5 * time.Second,
				WaitForCancellation: true,
			}
			disconnectedCtx = workflow.WithActivityOptions(disconnectedCtx, sitdownOptions)

			var recoveryResult string
			recoveryErr := workflow.ExecuteActivity(disconnectedCtx, ra.Sitdown, robotURL).Get(disconnectedCtx, &recoveryResult)

			if recoveryErr != nil {
				logger.Error("Failed to execute recovery (SitDown)", "error", recoveryErr)
				// 這裡可能需要根據情況決定回傳什麼錯誤，但通常原始的 Cancel 原因比較重要
				return "", err
			}

			return "", temporal.NewCanceledError("workflow canceled due to low battery")
		}

		return "", err
	}

	return result, nil
}
