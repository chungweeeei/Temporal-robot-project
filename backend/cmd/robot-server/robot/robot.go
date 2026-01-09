package robot

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type MockRobot struct {
	mu    sync.Mutex
	State RobotState
}

type RobotState struct {
	CurrentAction int
	BatteryLevel  int
	// [New] Position
	X        float64
	Y        float64
	IsMoving bool
}

func New() *MockRobot {
	r := &MockRobot{
		State: RobotState{
			CurrentAction: 2,
			BatteryLevel:  30,
		},
	}

	// Simulate battery drainage
	// go func() {
	// 	for {
	// 		r.mu.Lock()
	// 		if r.State.BatteryLevel > 1 {
	// 			r.State.BatteryLevel -= 1
	// 		}
	// 		r.mu.Unlock()

	// 		// Drain battery every 5 seconds
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	return r
}

// Handle Request
func (r *MockRobot) HandleRequest(request CallServiceRequest) CallServiceResponse {

	fmt.Println("Receive Service Request:", request)

	var args BaseRequestArgs
	if err := json.Unmarshal([]byte(request.Args.Data), &args); err != nil {
		fmt.Println("Error unmarshaling request args:", err)
		respData := BaseResponse{
			ApiID: args.ApiID,
			Status: StatusDetail{
				Code:    -1,
				Message: "Invalid request arguments",
			},
		}
		bytes, _ := json.Marshal(respData)
		return CallServiceResponse{
			Op:      "service_response",
			Service: request.Service,
			Values: struct {
				Data string `json:"data"`
			}{Data: string(bytes)},
			Result: false,
		}
	}

	switch args.ApiID {
	case 1009:
		r.mu.Lock()
		currentAction := r.State.CurrentAction
		currentBattery := r.State.BatteryLevel
		r.mu.Unlock()

		respData := RobotStateResponse{
			ApiID:         args.ApiID,
			CurrentAction: currentAction,
			BatteryLevel:  currentBattery,
			X:             float64(r.State.X),
			Y:             float64(r.State.Y),
			IsMoving:      r.State.IsMoving,
			Status: StatusDetail{
				Code:    0,
				Message: "Success",
			},
		}
		bytes, _ := json.Marshal(respData)
		return CallServiceResponse{
			Op:      "service_response",
			Service: request.Service,
			Values: struct {
				Data string `json:"data"`
			}{Data: string(bytes)},
			Result: true,
		}
	case 1012:
		var moveArgs struct {
			ApiID       int     `json:"api_id"`
			MissionID   string  `json:"mission_id"`
			X           float64 `json:"x"`
			Y           float64 `json:"y"`
			Orientation float64 `json:"orientation"`
		}
		json.Unmarshal([]byte(request.Args.Data), &moveArgs)

		r.mu.Lock()
		isClose := isClose(r.State.X, r.State.Y, moveArgs.X, moveArgs.Y, 0.1) // 假設 0.1 為誤差容忍值
		isMoving := r.State.IsMoving
		r.mu.Unlock()

		if isClose {
			fmt.Println("Robot is already at target location.")
			respData := BaseResponse{
				ApiID: args.ApiID,
				Status: StatusDetail{
					Code:    0,
					Message: "Already at target",
				},
			}
			bytes, _ := json.Marshal(respData)
			return CallServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{Data: string(bytes)},
				Result: true,
			}
		}

		if isMoving {
			fmt.Println("Robot is busy moving.")
			respData := BaseResponse{
				ApiID: args.ApiID,
				Status: StatusDetail{
					Code:    -1,
					Message: "Robot is busy",
				},
			}
			bytes, _ := json.Marshal(respData)
			return CallServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{Data: string(bytes)},
				Result: false,
			}
		}

		r.BackgroundMove(moveArgs.X, moveArgs.Y)

		respData := BaseResponse{
			ApiID: args.ApiID,
			Status: StatusDetail{
				Code:    0,
				Message: "Move command accepted (Async)",
			},
		}
		bytes, _ := json.Marshal(respData)

		return CallServiceResponse{
			Op:      "service_response",
			Service: request.Service,
			Values: struct {
				Data string `json:"data"`
			}{Data: string(bytes)},
			Result: true,
		}
	case 1013:
		var motionArgs struct {
			ApiID  int `json:"api_id"`
			Action int `json:"action"`
		}
		if err := json.Unmarshal([]byte(request.Args.Data), &motionArgs); err != nil {
			respData := BaseResponse{
				ApiID: args.ApiID,
				Status: StatusDetail{
					Code:    -1,
					Message: "Invalid motion control arguments",
				},
			}
			bytes, _ := json.Marshal(respData)
			return CallServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{Data: string(bytes)},
				Result: false,
			}
		}
		// Assume motion need 30 seconds to complete
		switch motionArgs.Action {
		case 3:
			time.Sleep(30 * time.Second)
		case 4:
			time.Sleep(5 * time.Second)
		}

		// Update robot state based on action
		r.mu.Lock()
		r.State.CurrentAction = motionArgs.Action
		r.mu.Unlock()

		respData := BaseResponse{
			ApiID: args.ApiID,
			Status: StatusDetail{
				Code:    0,
				Message: "Motion command executed",
			},
		}
		bytes, _ := json.Marshal(respData)
		return CallServiceResponse{
			Op:      "service_response",
			Service: request.Service,
			Values: struct {
				Data string `json:"data"`
			}{Data: string(bytes)},
			Result: true,
		}
	default:
		fmt.Println("Receive unknown api_id request")
		respData := BaseResponse{
			ApiID: args.ApiID,
			Status: StatusDetail{
				Code:    -1,
				Message: "Unknown ApiID",
			},
		}
		bytes, _ := json.Marshal(respData)
		return CallServiceResponse{
			Op:      "service_response",
			Service: request.Service,
			Values: struct {
				Data string `json:"data"`
			}{Data: string(bytes)},
			Result: false,
		}
	}
}

func (r *MockRobot) BackgroundMove(targetX, targetY float64) {

	go func() {
		fmt.Printf("Background move started: Target (%f, %f)\n", targetX, targetY)

		r.mu.Lock()
		r.State.IsMoving = true
		startX := r.State.X
		startY := r.State.Y
		r.mu.Unlock()

		// simple simulation: move step by step
		// assume it takes 10 seconds to reach the target
		// update position every 0.1 second
		const durationSeconds = 60
		const updatesPerSecond = 10
		steps := durationSeconds * updatesPerSecond

		interval := time.Duration(1000/updatesPerSecond) * time.Millisecond

		for i := 1; i < steps; i++ {
			time.Sleep(interval)

			r.mu.Lock()
			progress := float64(i) / float64(steps)

			r.State.X = startX + (targetX-startX)*progress
			r.State.Y = startY + (targetY-startY)*progress

			// 僅在每秒時輸出 Log，避免 Console 被刷爆
			if i%updatesPerSecond == 0 {
				fmt.Printf("Robot moving... (%.0f%%) Pos(%.2f, %.2f)\n", progress*100, r.State.X, r.State.Y)
			}
			r.mu.Unlock()
		}

		r.mu.Lock()
		r.State.X = targetX
		r.State.Y = targetY
		r.State.IsMoving = false
		r.State.CurrentAction = 2 // 回到 IDLE 狀態 (假設 2 是 IDLE)
		fmt.Println("Robot arrived at target.")
		r.mu.Unlock()
	}()

}

func isClose(x1, y1, x2, y2, tolerance float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	return (dx*dx + dy*dy) < (tolerance * tolerance)
}
