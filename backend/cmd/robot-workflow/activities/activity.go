package activities

type RobotActivities struct {
	Client *RobotClient
}

type ServiceRequest struct {
	Op      string      `json:"op"`
	Service string      `json:"service"`
	Type    string      `json:"type"`
	Args    RequestArgs `json:"args"`
}

type RequestArgs struct {
	Data string `json:"data"`
}

type ServiceResponse struct {
	Op       string         `json:"op"`
	Services string         `json:"service"`
	Values   ResponseValues `json:"values"`
	Result   bool           `json:"result"`
}

type ResponseValues struct {
	Data string `json:"data"`
}

func NewRobotActivities() *RobotActivities {
	return &RobotActivities{
		Client: NewRobotClient(),
	}
}
