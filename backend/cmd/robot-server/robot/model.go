package robot

import (
	"log"
	"os"
	"sync"
)

type RobotState struct {
	BatteryLevel int
	// [New]
	X           float64
	Y           float64
	Orientation float64
}
type MockRobot struct {
	Mu       sync.Mutex
	State    RobotState
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func New() *MockRobot {
	r := &MockRobot{
		Mu: sync.Mutex{},
		State: RobotState{
			BatteryLevel: 30,
			X:            0.0,
			Y:            0.0,
			Orientation:  0.0,
		},
		InfoLog:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return r
}
