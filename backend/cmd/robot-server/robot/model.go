package robot

import (
	"log"
	"os"
	"sync"
)

type MissionCode int

const (
	MissionCodeInit MissionCode = iota
	MissionCodeStart
	MissionSuccess
	MissionFailed
	MissionAbort
)

type RobotState struct {
	BatteryLevel int
	// [New]
	X           float64
	Y           float64
	Orientation float64
	// [New]
	MissionID string
	Mission   struct {
		Code    MissionCode
		Message string
	}
}
type MockRobot struct {
	Mu       sync.Mutex
	State    RobotState
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	StopChan chan bool
}

func New() *MockRobot {
	r := &MockRobot{
		Mu: sync.Mutex{},
		State: RobotState{
			BatteryLevel: 30,
			X:            0.0,
			Y:            0.0,
			Orientation:  0.0,
			MissionID:    "0",
			Mission: struct {
				Code    MissionCode
				Message string
			}{
				Code:    MissionCodeInit,
				Message: "INIT",
			},
		},
		InfoLog:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
		StopChan: make(chan bool),
	}
	return r
}
