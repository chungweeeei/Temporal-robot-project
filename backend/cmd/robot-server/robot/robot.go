package robot

import (
	"encoding/json"
	"fmt"
	"sync"
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
			BatteryLevel:  100,
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
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	return r
}

// Handle Request
func (r *MockRobot) HandleRequest(request CallServiceRequest) CallServiceResponse {

	r.mu.Lock()
	defer r.mu.Unlock()

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
		fmt.Println("Receive fetch robot status request")
		respData := RobotStateResponse{
			ApiID:         args.ApiID,
			CurrentAction: r.State.CurrentAction,
			BatteryLevel:  r.State.BatteryLevel,
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
		fmt.Println("Receive execute robot motion request")
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
		// Update robot state based on action
		r.State.CurrentAction = motionArgs.Action

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
