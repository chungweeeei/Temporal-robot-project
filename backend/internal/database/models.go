package database

import (
	"log"

	"github.com/chungweeeei/Temporal-robot-project/internal/repository/dao"
	"github.com/chungweeeei/Temporal-robot-project/internal/repository/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func New(dbPool *gorm.DB) Models {

	db = dbPool

	// Do auto migration
	err := db.AutoMigrate(&models.ActivityDefinition{}, &models.Workflow{})
	if err != nil {
		log.Println("Failed to auto migrate workflows table")
	}

	return Models{
		Workflow: dao.NewWorkflowDAO(db),
		Activity: dao.NewActivityDAO(db),
	}
}

type Models struct {
	Workflow models.WorkflowInterface
	Activity models.ActivityInterface
}
