package robot

type CallServiceRequest struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Type    string `json:"type"`
	Args    struct {
		Data string `json:"data"`
	} `json:"args"`
}

type CallServiceResponse struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Values  struct {
		Data string `json:"data"`
	} `json:"values"`
	Result bool `json:"result"`
}

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
	Status        StatusDetail `json:"status"`
}
