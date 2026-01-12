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

func (ra *RobotActivities) GetStatus(ctx context.Context) (RobotStatus, error) {

	return executeWithHeartbeat(ctx, func() (RobotStatus, error) {

		data := map[string]int{
			"api_id": RobotStatusID,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return RobotStatus{}, err
		}

		responseStr, err := ra.Client.CallService(ctx, string(dataBytes))
		if err != nil {
			return RobotStatus{}, err
		}

		var status RobotStatus
		err = json.Unmarshal([]byte(responseStr), &status)
		if err != nil {
			return RobotStatus{}, err
		}

		return status, nil
	})
}
