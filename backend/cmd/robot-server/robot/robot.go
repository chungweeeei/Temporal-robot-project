package robot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/pkg"
)

// Handle Request
func (r *MockRobot) HandleRequest(request pkg.ServiceRequest) pkg.ServiceResponse {

	r.InfoLog.Println("Receive Service Request:", request)

	switch request.Service {
	case "/api/system":
		// pre-processing data casting
		requestDataStr, ok := request.Args.Data.(string)
		if !ok {
			r.ErrorLog.Println("Invalid request data format")
			return pkg.ServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{
					Data: "",
				},
			}
		}
		requestDataBytes := []byte(requestDataStr)

		// First validate request payload (Partial Decode)
		var args BaseRequestArgs
		if err := json.Unmarshal(requestDataBytes, &args); err != nil {
			r.ErrorLog.Println("Error unmarshaling request args:", err)

			respData := BaseResponse{
				ApiID: args.ApiID,
				Status: StatusDetail{
					Code:    INVALID_INPUT,
					Message: "Failed to unmarshaling request arguments, please check your input",
				},
			}
			bytes, _ := json.Marshal(respData)
			return pkg.ServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{
					Data: string(bytes),
				},
			}
		}

		// Handle different service
		switch args.ApiID {
		case RobotMoveCommandID:
			return r.MoveToLocation(request.Service, requestDataBytes)
		case RobotMotionControlID:
			return r.HandleMotionControl(request.Service, requestDataBytes)
		default:
			return r.HandleUnknownRequest(args.ApiID, request.Service)
		}
	case "/set_angle_tag":
		// pre-processing data casting
		requestAngle, ok := request.Args.Data.(float64)
		if !ok {
			r.ErrorLog.Println("Invalid request data format")
			return pkg.ServiceResponse{
				Op:      "service_response",
				Service: request.Service,
				Values: struct {
					Data string `json:"data"`
				}{
					Data: "",
				},
			}
		}
		requestDataBytes := []byte(fmt.Sprintf("%f", requestAngle))
		return r.HandleHeadAngle(request.Service, requestDataBytes)
	default:
		return r.HandleUnknownService(request.Service)
	}
}

func (r *MockRobot) Move(targetX, targetY float64) {

	fmt.Printf("Background move started: Target (%f, %f)\n", targetX, targetY)

	r.Mu.Lock()
	startX := r.State.X
	startY := r.State.Y
	r.Mu.Unlock()

	// simple simulation: move step by step
	// assume it takes 10 seconds to reach the target
	// update position every 0.1 second
	const durationSeconds = 60
	const updatesPerSecond = 10
	steps := durationSeconds * updatesPerSecond

	interval := time.Duration(1000/updatesPerSecond) * time.Millisecond

	for i := 1; i < steps; i++ {
		time.Sleep(interval)

		r.Mu.Lock()
		progress := float64(i) / float64(steps)

		r.State.X = startX + (targetX-startX)*progress
		r.State.Y = startY + (targetY-startY)*progress

		if i%updatesPerSecond == 0 {
			fmt.Printf("Robot moving... (%.0f%%) Pos(%.2f, %.2f)\n", progress*100, r.State.X, r.State.Y)
		}
		r.Mu.Unlock()
	}

	r.Mu.Lock()
	r.State.X = targetX
	r.State.Y = targetY
	r.Mu.Unlock()

	r.InfoLog.Printf("Robot reached target location (%.2f, %.2f)\n", targetX, targetY)
}
