package activities

type RobotActivities struct {
	Client *RobotClient
}

func NewRobotActivities() *RobotActivities {
	return &RobotActivities{
		Client: NewRobotClient(),
	}
}
