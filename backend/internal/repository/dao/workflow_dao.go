package dao

import (
	"errors"
	"fmt"

	"github.com/chungweeeei/Temporal-robot-project/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkflowDAO struct {
	DB *gorm.DB
}

func NewWorkflowDAO(db *gorm.DB) *WorkflowDAO {
	return &WorkflowDAO{
		DB: db,
	}
}

func (dao *WorkflowDAO) Upsert(workflow models.Workflow) (string, error) {

	result := dao.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "workflow_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"workflow_name", "nodes", "updated_at"}),
	}).Create(&workflow)

	if result.Error != nil {
		return "", errors.New("failed to insert workflow")
	}

	return workflow.WorkflowID, nil
}

func (dao *WorkflowDAO) Get() ([]models.Workflow, error) {

	workflows := []models.Workflow{}
	result := dao.DB.Order("created_at").Find(&workflows)
	if result.Error != nil {
		return nil, errors.New("failed to retrieve workflows")
	}

	return workflows, nil
}

func (dao *WorkflowDAO) GetByID(id string) (*models.Workflow, error) {

	workflow := models.Workflow{}
	result := dao.DB.Where("workflow_id = ?", id).First(&workflow)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workflow with %s not found", id)
		}
		fmt.Println(result.Error)
		return nil, errors.New("failed to retrieve workflow by id")
	}

	return &workflow, nil
}

func (dao *WorkflowDAO) Delete(workflowId string) error {

	result := dao.DB.Where("workflow_id = ?", workflowId).Delete(&models.Workflow{})
	if result.Error != nil {
		return errors.New("failed to delete workflow")
	}
	return nil
}
