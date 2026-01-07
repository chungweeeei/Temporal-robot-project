package main

import (
	"log"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "ROBOT_TASK_QUEUE", worker.Options{})

	activities := activities.NewRobotActivities()
	w.RegisterWorkflow(workflows.RobotActionWorkflow)
	w.RegisterWorkflow(workflows.HealthCheckWorkflow)
	w.RegisterActivity(activities)

	// automatically trigger health check workflow
	// start a new goroutine to start the health check workflow, so that it doesn't block the worker
	// go func(){
	// 	ctx := context.Background()

	// }

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
