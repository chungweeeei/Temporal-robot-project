package models

import (
	"time"

	"gorm.io/datatypes"
)

type ActivityDefinition struct {
	ID           int            `json:"id" gorm:"primaryKey autoIncrement"`
	Name         string         `json:"name" gorm:"not null"`
	ActivityType string         `json:"activity_type" gorm:"uniqueIndex;not null"`
	NodeType     string         `json:"node_type" gorm:"default:'default'"`
	InputSchema  datatypes.JSON `json:"input_schema" gorm:"type:json"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}
