package data

import (
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

func New(dbPool *gorm.DB) Models {

	db = dbPool

	// Do auto migration
	err := db.AutoMigrate(&Workflow{})
	if err != nil {
		log.Println("Failed to auto migrate workflows table")
	}

	return Models{
		Workflow: &Workflow{},
	}
}

type Models struct {
	Workflow WorkflowInterface
}
