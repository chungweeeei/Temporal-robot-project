package robot

type BaseRequestArgs struct {
	ApiID int `json:"api_id"`
}

type BaseResponse struct {
	ApiID  int          `json:"api_id"`
	Status StatusDetail `json:"status"`
}

type StatusDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RobotStateResponse struct {
	ApiID         int          `json:"api_id"`
	CurrentAction int          `json:"current_action"`
	BatteryLevel  int          `json:"battery_level"`
	X             float64      `json:"x"`
	Y             float64      `json:"y"`
	IsMoving      bool         `json:"is_moving"`
	Status        StatusDetail `json:"status"`
}
