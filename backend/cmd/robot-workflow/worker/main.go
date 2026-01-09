package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	// Register low-level slog hanlder setting level to INFO
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Register Temporal client
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		Logger:   logger,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "ROBOT_TASK_QUEUE", worker.Options{})

	activities := activities.NewRobotActivities()
	w.RegisterWorkflow(workflows.RobotActionWorkflow)
	w.RegisterWorkflow(workflows.RobotMonitorWorkflow)
	w.RegisterActivity(activities)

	// // Automatically start the monitor Workflow
	// go func() {
	// 	workflowOptions := client.StartWorkflowOptions{
	// 		ID:        "robot_monitor_workflow",
	// 		TaskQueue: "ROBOT_TASK_QUEUE",
	// 	}
	// 	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RobotMonitorWorkflow)
	// 	if err != nil {
	// 		log.Printf("Unable to execute workflow (might be already running): %v", err)
	// 	} else {
	// 		log.Printf("Started workflow WorkflowID: %s RunID: %s", we.GetID(), we.GetRunID())
	// 	}
	// }()

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
