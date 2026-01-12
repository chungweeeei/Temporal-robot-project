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

	if payload.RootNodeID == "" {
		return "", fmt.Errorf("rootNodeId is missing")
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	currentStep := "Initializing"
	workflow.SetQueryHandler(ctx, "get_step", func() (string, error) {
		return currentStep, nil
	})

	currentNodeID := payload.RootNodeID
	for {
		currentNode, exists := payload.Nodes[currentNodeID]
		if !exists {
			return "", fmt.Errorf("node with ID %s not found", currentNodeID)
		}

		currentStep = fmt.Sprintf("Executing node %s of type %s", currentNodeID, currentNode.Type)
		switch currentNode.Type {
		case pkg.ActivityStandUp, pkg.ActivitySitDown, pkg.ActivityHead, pkg.ActivityMove, pkg.ActivityTTS:
			// Execute robot activity
			var result string
			future := workflow.ExecuteActivity(ctx, string(currentNode.Type), currentNode.Params)
			err := future.Get(ctx, &result)
			if err != nil {
				logger.Error("Activity failed", "error", err)
			}
			// Determine next node based on success or failure
			if err != nil {
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
				return "", fmt.Errorf("invalid or missing duration parameter for sleep activity")
			}
			duration := int(durationFloat)
			workflow.Sleep(ctx, time.Millisecond*time.Duration(duration))
			// Move to next node
			if currentNode.Transitions.Next == "" {
				logger.Info("Workflow completed successfully")
				return "Workflow completed successfully", nil
			}
			currentNodeID = currentNode.Transitions.Next
		case pkg.ActivityEnd:
			logger.Info("Workflow reached end node")
			return "Workflow completed successfully", nil
		case pkg.ActivityStart:
			currentNodeID = currentNode.Transitions.Next
		default:
			return "", fmt.Errorf("unsupported activity type: %s", currentNode.Type)
		}
	}
}
