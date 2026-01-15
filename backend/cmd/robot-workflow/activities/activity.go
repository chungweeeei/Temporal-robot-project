package activities

type RobotActivities struct {
	Client      *RobotClient
	StatusCache *StatusCache
}

func NewRobotActivities(robotIP string, statusCache *StatusCache) *RobotActivities {
	return &RobotActivities{
		Client:      NewRobotClient(robotIP),
		StatusCache: statusCache,
	}
}
