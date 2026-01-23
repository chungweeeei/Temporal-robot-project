package models

type WorkflowInterface interface {
	Upsert(workflow Workflow) (string, error)
	Get() ([]Workflow, error)
	GetByID(id string) (*Workflow, error)
	Delete(workflowId string) error
}

type ActivityInterface interface {
	Get() ([]ActivityDefinition, error)
}
