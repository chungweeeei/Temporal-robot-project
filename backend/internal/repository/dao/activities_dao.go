package dao

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/chungweeeei/Temporal-robot-project/internal/repository/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivityDAO struct {
	DB *gorm.DB
}

func NewActivityDAO(db *gorm.DB) *ActivityDAO {

	seeds := getActivitySeeds()

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "activity_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "node_type", "input_schema", "updated_at"}),
	}).Create(&seeds).Error

	if err != nil {
		log.Printf("Failed to seed activities: %v", err)
	} else {
		log.Println("Activities definitions synced (Upsert).")
	}

	// seed Activity Definitions
	return &ActivityDAO{
		DB: db,
	}
}

func getActivitySeeds() []models.ActivityDefinition {

	// Define activities definitions
	moveSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"x":           map[string]interface{}{"type": "number", "title": "X (m)", "step": 0.1, "default": 0.0},
			"y":           map[string]interface{}{"type": "number", "title": "Y (m)", "step": 0.1, "default": 0.0},
			"orientation": map[string]interface{}{"type": "number", "title": "Orientation (Degree)", "step": 0.1, "default": 0.0},
		},
		"required": []string{"x", "y"},
	}
	moveJSON, _ := json.Marshal(moveSchema)

	sleepSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"duration": map[string]interface{}{"type": "number", "title": "Duration (milliseconds)", "default": 3000},
		},
		"required": []string{"duration"},
	}
	sleepJSON, _ := json.Marshal(sleepSchema)

	ttsSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"text": map[string]interface{}{"type": "string", "title": "Message"},
		},
		"required": []string{"text"},
	}
	ttsJSON, _ := json.Marshal(ttsSchema)

	headSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"angle": map[string]interface{}{"type": "number", "title": "Angle (degree)", "default": 0.0},
		},
		"required": []string{"angle"},
	}
	headJSON, _ := json.Marshal(headSchema)

	return []models.ActivityDefinition{
		{
			Name:         "Move",
			ActivityType: "Move",
			NodeType:     "action",
			InputSchema:  datatypes.JSON(moveJSON),
		},
		{
			Name:         "Sleep",
			ActivityType: "Sleep",
			NodeType:     "action",
			InputSchema:  datatypes.JSON(sleepJSON),
		},
		{
			Name:         "Standup",
			ActivityType: "Standup",
			NodeType:     "action",
			InputSchema:  datatypes.JSON{},
		},
		{
			Name:         "Sitdown",
			ActivityType: "Sitdown",
			NodeType:     "action",
			InputSchema:  datatypes.JSON{},
		},
		{
			Name:         "TTS",
			ActivityType: "TTS",
			NodeType:     "action",
			InputSchema:  datatypes.JSON(ttsJSON),
		},
		{
			Name:         "Head",
			ActivityType: "Head",
			NodeType:     "action",
			InputSchema:  datatypes.JSON(headJSON),
		},
	}

}

func (dao *ActivityDAO) Get() ([]models.ActivityDefinition, error) {

	activities := []models.ActivityDefinition{}
	result := dao.DB.Find(&activities)
	if result.Error != nil {
		return nil, errors.New("failed to retrieve activities")
	}

	return activities, nil
}
