package data

import (
	"errors"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Workflow struct {
	WorkflowID   string         `json:"workflow_id" gorm:"primaryKey"`
	WorkflowName string         `json:"workflow_name" gorm:"unique; not null; VARCHAR(255)"`
	RootNodeID   string         `json:"root_node_id" gorm:"not null; VARCHAR(255) default:'start'"`
	Nodes        datatypes.JSON `json:"nodes" gorm:"type:json; not null"`
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime"`
}

func (w *Workflow) Upsert(workflow Workflow) (string, error) {

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "workflow_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"workflow_name", "nodes", "updated_at"}),
	}).Create(&workflow)

	if result.Error != nil {
		return "", errors.New("failed to insert workflow")
	}

	return workflow.WorkflowID, nil
}

func (w *Workflow) Get() ([]Workflow, error) {

	workflows := []Workflow{}
	result := db.Find(&workflows).Order("created_at")
	if result.Error != nil {
		return nil, errors.New("failed to retrieve workflows")
	}

	return workflows, nil
}

func (w *Workflow) GetByID(id string) (*Workflow, error) {

	workflow := Workflow{}
	result := db.Where("workflow_id = ?", id).First(&workflow)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow with %s not found", id)
		}
		fmt.Println(result.Error)
		return nil, errors.New("failed to retrieve workflow by id")
	}

	return &workflow, nil
}
