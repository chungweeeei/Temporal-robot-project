package activity

type RobotActivities struct {
	Client      *RobotClient
	CacheStatus *CacheStatus
}

func NewRobotActivities(robotIP string, cacheStatus *CacheStatus) *RobotActivities {
	return &RobotActivities{
		Client:      NewRobotClient(robotIP),
		CacheStatus: cacheStatus,
	}
}
