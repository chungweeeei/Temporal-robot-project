package activities

type RobotActivities struct {
	Client      *RobotClient
	StatusCache *StatusCache
}

func NewRobotActivities(statusCache *StatusCache) *RobotActivities {
	return &RobotActivities{
		Client:      NewRobotClient(),
		StatusCache: statusCache,
	}
}
