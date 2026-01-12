package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func (app *Config) triggerWorkflowHandler(c *gin.Context) {

	var payload pkg.WorkflowPayload
	// 1. bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		app.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	// 2. make sure workflow ID is set
	if payload.WorkflowID == "" {
		payload.WorkflowID = "robot-routine-" + uuid.New().String()
	}

	app.InfoLog.Println("Received workflow trigger request:", payload)

	// 3. setting workflow options
	workflowOptions := client.StartWorkflowOptions{
		ID:        payload.WorkflowID,
		TaskQueue: "ROBOT_TASK_QUEUE",
	}

	we, err := app.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RobotWorkflow, payload)
	if err != nil {
		app.ErrorLog.Println("Unable to start workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to start workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow started successfully",
		"workflow_id": we.GetID(),
	})
}
