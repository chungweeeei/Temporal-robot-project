package main

import (
	"context"
	"log"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	workflowID := "robot_action_workflow_001"
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "ROBOT_TASK_QUEUE",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.RobotActionWorkflow, workflows.RobotWorkflowRequest{Action: workflows.Move, MoveTarget: &activities.MoveRequest{X: 10.0, Y: 10.0, Orientation: 0.0}})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}

	log.Println("Workflow result:", result)
}
