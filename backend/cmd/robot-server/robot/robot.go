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
