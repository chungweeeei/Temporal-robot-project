package workflows

import (
	"fmt"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func RobotWorkflow(ctx workflow.Context, payload pkg.WorkflowPayload) (string, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("Robot workflow started", "payload", payload)

	// validate received payload
	if payload.RootNodeID == "" {
		return "", fmt.Errorf("rootNodeId is missing")
	}

	// prepared activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Register signal for stop & resume workflow
	pause := false
	var cancelCurrentActivity func()
	signalChan := workflow.GetSignalChannel(ctx, "control-signal")

	currentStep := "Initializing"

	// Background listener for control signal
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			var signal string
			signalChan.Receive(ctx, &signal)
			logger.Info("Received control signal", "signal", signal)
			switch signal {
			case "pause":
				pause = true
				if cancelCurrentActivity != nil {
					cancelCurrentActivity()
				}
			case "resume":
				pause = false
			}
		}
	})

	workflow.SetQueryHandler(ctx, "get_step", func() (string, error) {
		if pause {
			return "Paused", nil
		}
		return currentStep, nil
	})

	currentNodeID := payload.RootNodeID
	for {
		// 使用 workflow.Await 來等待 paused 狀態解除
		// 這裡會阻塞直到匿名函數返回 true (即 !paused)
		// 這樣在任何 Activity 執行"前"，都會檢查是否暫停
		workflow.Await(ctx, func() bool { return !pause })

		// Register children cancel context
		childCtx, cancel := workflow.WithCancel(ctx)
		cancelCurrentActivity = cancel

		currentNode, exists := payload.Nodes[currentNodeID]
		if !exists {
			cancel()
			return "", fmt.Errorf("node with ID %s not found", currentNodeID)
		}

		currentStep = string(currentNode.Type)
		switch currentNode.Type {
		case pkg.ActivityStandUp, pkg.ActivitySitDown, pkg.ActivityHead, pkg.ActivityMove, pkg.ActivityTTS:
			// Execute robot activity
			var result string
			err := workflow.ExecuteActivity(childCtx, string(currentNode.Type), currentNode.Params).Get(childCtx, &result)

			// clean up cancle function
			cancelCurrentActivity = nil
			cancel()

			if temporal.IsCanceledError(err) {
				logger.Info("Activity was cancelled due to pause signal", "activity", string(currentNode.Type))
				currentStep = "Paused"
				continue
			}

			// Determine next node based on success or failure
			if err != nil {
				logger.Error("Activity failed", "error", err)
				if currentNode.Transitions.Failure == "" {
					return "", fmt.Errorf("no failure transition defined for node %s", currentNodeID)
				}
				currentNodeID = currentNode.Transitions.Failure
			} else {
				if currentNode.Transitions.Next == "" {
					logger.Info("Workflow completed successfully")
					return "Workflow completed successfully", nil
				}
				currentNodeID = currentNode.Transitions.Next
			}

		case pkg.ActivitySleep:
			// Sleep activity
			durationFloat, ok := currentNode.Params["duration"].(float64)
			if !ok {
				cancel()
				return "", fmt.Errorf("invalid or missing duration parameter for sleep activity")
			}
			duration := int(durationFloat)
			err := workflow.Sleep(ctx, time.Millisecond*time.Duration(duration))

			cancelCurrentActivity = nil
			cancel()

			if temporal.IsCanceledError(err) {
				logger.Info("Sleep activity was cancelled due to pause signal")
				continue
			}

			// Move to next node
			if currentNode.Transitions.Next == "" {
				logger.Info("Workflow completed successfully")
				return "Workflow completed successfully", nil
			}
			currentNodeID = currentNode.Transitions.Next
		case pkg.ActivityEnd:
			cancel()
			logger.Info("Workflow reached end node")
			return "Workflow completed successfully", nil
		case pkg.ActivityStart:
			cancel()
			currentNodeID = currentNode.Transitions.Next
		default:
			cancel()
			return "", fmt.Errorf("unsupported activity type: %s", currentNode.Type)
		}
	}
}
