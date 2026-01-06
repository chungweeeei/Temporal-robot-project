package main

import (
	"log"

	"github.com/chungweeeei/Temporal-robot-project/temporal_worker/robot_routine/workflows"
	sharedactivities "github.com/chungweeeei/Temporal-robot-project/temporal_worker/shared_activities"
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

	w := worker.New(c, "ROBOT_ROUTINE_TASK_QUEUE", worker.Options{})

	activities := sharedactivities.NewRobotActivities()

	w.RegisterWorkflow(workflows.RobotRoutine)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("Unable to start Worker", err)
	}
}
