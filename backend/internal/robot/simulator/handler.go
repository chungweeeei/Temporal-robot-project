package simulator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/pkg"
)

func (r *MockRobot) MoveToLocation(service string, request []byte) pkg.ServiceResponse {

	var moveArgs struct {
		ApiID       int     `json:"api_id"`
		MissionID   string  `json:"mission_id"`
		X           float64 `json:"x"`
		Y           float64 `json:"y"`
		Orientation float64 `json:"orientation"`
	}
	err := json.Unmarshal(request, &moveArgs)
	if err != nil {
		r.ErrorLog.Println("Error unmarshaling move args:", err)
		respData := BaseResponse{
			ApiID: RobotMoveCommandID,
			Status: StatusDetail{
				Code:    INVALID_INPUT,
				Message: "Failed to unmarshaling move arguments, please check your input",
			},
		}
		bytes, _ := json.Marshal(respData)
		return pkg.ServiceResponse{
			Op:      "service_response",
			Service: service,
			Values: struct {
				Data string `json:"data"`
			}{
				Data: string(bytes),
			},
		}
	}

	// starting goroutine move to target point
	go r.Move(moveArgs.MissionID, moveArgs.X, moveArgs.Y)

	respData := BaseResponse{
		ApiID: RobotMoveCommandID,
		Status: StatusDetail{
			Code:    SUCCESS,
			Message: "Move command accepted (Async)",
		},
	}
	bytes, _ := json.Marshal(respData)

	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleStopCommand(service string, request []byte) pkg.ServiceResponse {

	// send stop signal to robot
	r.StopChan <- true

	// return success response
	respData := BaseResponse{
		ApiID: RobotStopActionID,
		Status: StatusDetail{
			Code:    SUCCESS,
			Message: "Stop command accepted",
		},
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleMotionControl(service string, request []byte) pkg.ServiceResponse {

	var motionArgs struct {
		ApiID  int `json:"api_id"`
		Action int `json:"action"`
	}

	err := json.Unmarshal(request, &motionArgs)
	if err != nil {
		respData := BaseResponse{
			ApiID: RobotMotionControlID,
			Status: StatusDetail{
				Code:    INVALID_INPUT,
				Message: "Failed to unmarshaling motion control arguments, please check your input",
			},
		}
		bytes, _ := json.Marshal(respData)
		return pkg.ServiceResponse{
			Op:      "service_response",
			Service: service,
			Values: struct {
				Data string `json:"data"`
			}{
				Data: string(bytes),
			},
		}
	}

	// mock 2 seconds delay for motion control execution
	time.Sleep(2 * time.Second)

	respData := BaseResponse{
		ApiID: RobotMotionControlID,
		Status: StatusDetail{
			Code:    SUCCESS,
			Message: "Motion command accepted",
		},
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleTTSCommand(service string, request []byte) pkg.ServiceResponse {

	var ttsArgs struct {
		ApiID     int     `json:"api_id"`
		Text      string  `json:"text"`
		VoiceName string  `json:"voice_name"`
		Speed     float64 `json:"speed"`
		Volume    float64 `json:"volume"`
	}

	if err := json.Unmarshal(request, &ttsArgs); err != nil {
		respData := BaseResponse{
			ApiID: RobotTTSCommandID,
			Status: StatusDetail{
				Code:    INVALID_INPUT,
				Message: "Failed to unmarshaling TTS command arguments, please check your input",
			},
		}
		bytes, _ := json.Marshal(respData)
		return pkg.ServiceResponse{
			Op:      "service_response",
			Service: service,
			Values: struct {
				Data string `json:"data"`
			}{
				Data: string(bytes),
			},
		}
	}

	// mock 2 seconds delay for motion control execution
	time.Sleep(2 * time.Second)

	respData := BaseResponse{
		ApiID: RobotTTSCommandID,
		Status: StatusDetail{
			Code:    SUCCESS,
			Message: "TTS command accepted",
		},
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleHeadAngle(service string, request []byte) pkg.ServiceResponse {

	// mock dealing with head angle setting
	time.Sleep(2 * time.Second)

	respData := StatusDetail{
		Code:    SUCCESS,
		Message: "Set Head angle accepted",
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleUnknownRequest(unknownId int, service string) pkg.ServiceResponse {
	r.ErrorLog.Println("Received Unknown API ID:", unknownId)
	respData := BaseResponse{
		ApiID: unknownId,
		Status: StatusDetail{
			Code:    NOT_EXIST,
			Message: fmt.Sprintf("Received Unknown API ID: %d", unknownId),
		},
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) HandleUnknownService(service string) pkg.ServiceResponse {
	r.ErrorLog.Println("Received Unknown Service:", service)
	respData := BaseResponse{
		ApiID: 0,
		Status: StatusDetail{
			Code:    NOT_EXIST,
			Message: fmt.Sprintf("Received Unknown Service: %s", service),
		},
	}
	bytes, _ := json.Marshal(respData)
	return pkg.ServiceResponse{
		Op:      "service_response",
		Service: service,
		Values: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}

func (r *MockRobot) GetRobotStatus() pkg.TopicResponse {

	qx, qy, qz, qw := transferOrientationToQuaternion(r.State.Orientation)

	robotStatus := RobotStatus{
		ApiID:        0,
		BatteryLevel: r.State.BatteryLevel,
		Pose: struct {
			Orientation struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
				W float64 `json:"w"`
			} `json:"orientation"`
			Position struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"position"`
		}{
			Orientation: struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
				W float64 `json:"w"`
			}{
				X: qx,
				Y: qy,
				Z: qz,
				W: qw,
			},
			Position: struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			}{
				X: r.State.X,
				Y: r.State.Y,
				Z: 0.0,
			},
		},
		MissionID: r.State.MissionID,
		Mission: struct {
			Code    MissionCode `json:"code"`
			Message string      `json:"message"`
		}{
			Code:    r.State.Mission.Code,
			Message: r.State.Mission.Message,
		},
		Timestamp: time.Now().Format("2006-01-02T15:04:05"),
	}

	respData := struct {
		DeviceName   string      `json:"device_name"`
		DeviceStatus RobotStatus `json:"device_status"`
	}{
		DeviceName:   "MockRobot",
		DeviceStatus: robotStatus,
	}

	bytes, _ := json.Marshal(respData)
	return pkg.TopicResponse{
		Op:    "publish",
		Topic: "/api/info",
		Msg: struct {
			Data string `json:"data"`
		}{
			Data: string(bytes),
		},
	}
}
