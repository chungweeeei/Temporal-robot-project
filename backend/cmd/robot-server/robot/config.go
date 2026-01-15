package robot

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
	RobotStatusID        = 1009
	RobotMoveCommandID   = 1012
	RobotMotionControlID = 1013
	RobotTTSCommandID    = 1014
	// TODO: define stop action ID temporally for testing pause/resume feature
	RobotStopActionID = 5000
)

const (
	StandUpActionID = 3
	SitDownActionID = 4
)
