package models

import "gorm.io/datatypes"

type ActivityDefinition struct {
	ID           int            `json:"id" gorm:"primaryKey autoIncrement"`
	Name         string         `json:"name" gorm:"not null"`
	ActivityType string         `json:"activity_type" gorm:"not null"`
	NodeType     string         `json:"node_type" gorm:"default:'default'"`
	InputSchema  datatypes.JSON `json:"input_schema" gorm:"type:json"`
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime"`
}
