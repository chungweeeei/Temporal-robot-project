package pkg

type ServiceRequest struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Type    string `json:"type"`
	Args    struct {
		Data any `json:"data"`
	} `json:"args"`
}

type ServiceResponse struct {
	Op      string `json:"op"`
	Service string `json:"service"`
	Values  struct {
		Data string `json:"data"`
	} `json:"values"`
}

type TopicRequest struct {
	Op           string `json:"op"`
	Topic        string `json:"topic"`
	Type         string `json:"type"`
	ThrottleRate int    `json:"throttle_rate"`
	QueueLength  int    `json:"queue_length"`
}

type TopicResponse struct {
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

type RetryPolicy struct {
	MaxAttempts        int32   `json:"maxAttempts"`
	InitialInterval    int32   `json:"initialInterval"` // ms
	BackoffCoefficient float64 `json:"backoffCoefficient,omitempty"`
	MaximumInterval    int32   `json:"maximumInterval,omitempty"` // ms
}

// WorkflowTransitions 定義狀態轉移
// 對應前端: next, failure 都可能是下一個節點 ID
type WorkflowTransitions struct {
	Next    string `json:"next,omitempty"`
	Failure string `json:"failure,omitempty"`
}

type WorkflowNode struct {
	ID          string                 `json:"id"`
	Type        ActivityType           `json:"type"`
	Params      map[string]interface{} `json:"params"`
	RetryPolicy *RetryPolicy           `json:"retryPolicy,omitempty"`
	Transitions WorkflowTransitions    `json:"transitions"`
}

type WorkflowPayload struct {
	WorkflowID string                  `json:"workflow_id,omitempty"`
	RootNodeID string                  `json:"root_node_id,omitempty"`
	Nodes      map[string]WorkflowNode `json:"nodes"`
}
