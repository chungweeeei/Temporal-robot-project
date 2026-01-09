package activities

import (
	"context"
	"encoding/json"
)

type RobotStatus struct {
	ApiID         int     `json:"api_id"`
	CurrentAction int     `json:"current_action"`
	BatteryLevel  int     `json:"battery_level"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	IsMoving      bool    `json:"is_moving"`
	Status        struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

func (ra *RobotActivities) GetStatus(ctx context.Context, url string) (RobotStatus, error) {

	errorChan := make(chan error)
	responseChan := make(chan RobotStatus)

	go func() {
		data := map[string]int{
			"api_id": RobotStatusID,
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			errorChan <- err
			return
		}

		response, err := ra.Client.CallService(ctx, url, string(dataBytes))
		if err != nil {
			errorChan <- err
			return
		}
		var status RobotStatus
		err = json.Unmarshal([]byte(response), &status)
		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- status
	}()

	select {
	case err := <-errorChan:
		return RobotStatus{}, err
	case status := <-responseChan:
		return status, nil
	case <-ctx.Done():
		return RobotStatus{}, ctx.Err()
	}

}
