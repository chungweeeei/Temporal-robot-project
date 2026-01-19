package simulator

type ErrorCode int

const (
	SUCCESS ErrorCode = iota
	FAILED
	INVALID_INPUT
	NOT_EXIST
	PERMISSION_DENIED
	ABORTED
	TYPE_ERROR
	REJECT
)

type ActionID int

const (
	RobotMoveCommandID   = 1005
	RobotStatusID        = 1009
	RobotMotionControlID = 1013
	RobotTTSCommandID    = 1014

	// TODO: define stop action ID temporally for testing pause/resume feature
	RobotStopActionID = 5000
)

const (
	StandUp   = 3
	StandDown = 4
	SitDown   = 6
)

type MissionCode int

const (
	MissionCodeInit MissionCode = iota
	MissionCodeStart
	MissionSuccess
	MissionFailed
	MissionAbort
)
