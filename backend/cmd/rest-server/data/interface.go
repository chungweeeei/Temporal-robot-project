package data

type WorkflowInterface interface {
	Upsert(workflow Workflow) (string, error)
	Get() ([]Workflow, error)
	GetByID(id string) (*Workflow, error)
}
