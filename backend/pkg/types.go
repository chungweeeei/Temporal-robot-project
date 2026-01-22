package pkg

type RobotServiceRequest struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Type    string `json:"type"`
	Args    struct {
		Data any `json:"data"`
	} `json:"args"`
}

type RobotServiceResponse struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Values  struct {
		Data string `json:"data"`
	} `json:"values"`
}

type RobotTopicRequest struct {
	Op           string `json:"op"`
	Topic        string `json:"topic"`
	Type         string `json:"type"`
	ThrottleRate int    `json:"throttle_rate"`
	QueueLength  int    `json:"queue_length"`
}

type RobotTopicResponse struct {
	Op    string `json:"op"`
	Topic string `json:"topic"`
	Msg   struct {
		Data string `json:"data"`
	} `json:"msg"`
}

type ActivityType string

const (
	ActivityStandUp ActivityType = "Standup"
	ActivitySitDown ActivityType = "Sitdown"
	ActivityHead    ActivityType = "Head"
	ActivityMove    ActivityType = "Move"
	ActivityTTS     ActivityType = "TTS"
	ActivitySleep   ActivityType = "Sleep"
	ActivityStart   ActivityType = "Start"
	ActivityEnd     ActivityType = "End"
)

type WorkflowTransitions struct {
	Next    string `json:"next,omitempty"`
	Failure string `json:"failure,omitempty"`
}

type WorkflowNode struct {
	ID          string                 `json:"id"`
	Type        ActivityType           `json:"type"`
	Params      map[string]interface{} `json:"params"`
	Transitions WorkflowTransitions    `json:"transitions"`
}

type WorkflowPayload struct {
	WorkflowID string                  `json:"workflow_id,omitempty"`
	RootNodeID string                  `json:"root_node_id,omitempty"`
	Nodes      map[string]WorkflowNode `json:"nodes"`
}
