package simulator

type BaseRequestArgs struct {
	ApiID int `json:"api_id"`
}

type BaseResponse struct {
	ApiID  int          `json:"api_id"`
	Status StatusDetail `json:"status"`
}

type StatusDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type RobotStatus struct {
	ApiID        int `json:"api_id"`
	BatteryLevel int `json:"battery_level"`
	Pose         struct {
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
	} `json:"pose"`
	MissionID string `json:"mission_id"`
	Mission   struct {
		Code    MissionCode `json:"code"`
		Message string      `json:"message"`
	} `json:"mission"`
}

type RobotStatusResponse struct {
	DeviceName   string `json:"device_name"`
	DevcieStatus string `json:"device_status"`
	Timestamp    string `json:"timestamp"`
}

type RobotDevice struct {
	Name string `json:"name"`
}

type DevicesInfo struct {
	Robots []RobotDevice `json:"robots"`
}

type FirmwareInfo struct {
	FwVersion string `json:"fw_version"`
	HwVersion string `json:"hw_version"`
	ModelName string `json:"model_name"`
}
type DeviceInfoResponse struct {
	ApiID    int          `json:"api_id"`
	Devices  DevicesInfo  `json:"devices"`
	Firmware FirmwareInfo `json:"firmware"`
	Status   StatusDetail `json:"status"`
}
