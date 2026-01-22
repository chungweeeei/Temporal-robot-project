package config

type MissionCode int

const (
	MissionInit MissionCode = iota
	MissionStart
	MissionSuccess
	MissionFailed
	MissionAbort
)

const (
	RobotMoveCommandID   = 1005
	RobotStatusID        = 1009
	RobotMotionControlID = 1013
	RobotTTSCommandID    = 1014
)

const (
	StandUpActionID = 3
	SitDownActionID = 4
)
