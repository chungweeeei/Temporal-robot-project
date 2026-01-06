package main

import (
	"log"

	actionworkflows "github.com/chungweeeei/robot/robot_action/workflows"
	routineworkflows "github.com/chungweeeei/robot/robot_routine/workflows"
	sharedactivities "github.com/chungweeeei/robot/shared_activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "ROBOT_TASK_QUEUE", worker.Options{})

	activities := sharedactivities.NewRobotActivities()
	w.RegisterWorkflow(actionworkflows.RobotAction)
	w.RegisterWorkflow(routineworkflows.RobotRoutine)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("Unable to start worker", err)
	}
}
