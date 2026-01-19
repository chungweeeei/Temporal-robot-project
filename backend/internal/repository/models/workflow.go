package models

import "gorm.io/datatypes"

type Workflow struct {
	WorkflowID   string         `json:"workflow_id" gorm:"primaryKey"`
	WorkflowName string         `json:"workflow_name" gorm:"unique; not null; VARCHAR(255)"`
	RootNodeID   string         `json:"root_node_id" gorm:"not null; VARCHAR(255) default:'start'"`
	Nodes        datatypes.JSON `json:"nodes" gorm:"type:json; not null"`
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime"`
}
