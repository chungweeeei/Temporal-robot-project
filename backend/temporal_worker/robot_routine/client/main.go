package main

import (
	"context"
	"log"
	"time"

	routineworkflows "github.com/chungweeeei/Temporal-robot-project/temporal_worker/robot_routine/workflows"
	"go.temporal.io/sdk/client"
)

func main() {

	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := "robot_routine_workflow_001"
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "ROBOT_TASK_QUEUE",
	}

	steps := []routineworkflows.RobotStep{
		{
			Action: routineworkflows.StandUp,
		},
		{
			Action:   routineworkflows.Sleep,
			Duration: 3 * time.Second, // 暫停 3 秒
		},
		{
			Action: routineworkflows.SitDown,
		},
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, routineworkflows.RobotRoutine, steps)
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
