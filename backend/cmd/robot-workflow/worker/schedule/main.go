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
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "ROBOT_SCHEDULE_QUEUE", worker.Options{})

	w.RegisterWorkflow(workflows.RobotScheduleWorkflow)
	w.RegisterActivity(activities.CheckTaskUpdates)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
