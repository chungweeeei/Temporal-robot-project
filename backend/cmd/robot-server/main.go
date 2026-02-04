package main

import (
	"log"
	"net/http"

	"github.com/chungweeeei/Temporal-robot-project/internal/robot/client"
	"github.com/chungweeeei/Temporal-robot-project/internal/robot/simulator"
)

func main() {
	robotSim := simulator.NewMockRobot()
	robotHandler := client.NewRobotHandler(robotSim)

	http.HandleFunc("/", robotHandler.HandleWS)
	log.Println("Mock Robot Server started on :9090")
	http.ListenAndServe(":9090", nil)
}
