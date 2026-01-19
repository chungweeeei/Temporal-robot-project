package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chungweeeei/Temporal-robot-project/cmd/rest-server/data"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SaveWorkflowRequest struct {
	WorkflowID   string                 `json:"workflow_id" binding:"required"`
	WorkflowName string                 `json:"workflow_name" binding:"required"`
	Nodes        map[string]interface{} `json:"nodes"`
}

func (app *Config) saveWorkflow(c *gin.Context) {

	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		app.ErrorLog.Println("Unable to marshal nodes:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to marshal nodes"})
		return
	}

	workflow := data.Workflow{
		WorkflowID:   req.WorkflowID,
		WorkflowName: req.WorkflowName,
		RootNodeID:   "start",
		Nodes:        nodes,
	}

	id, err := app.Model.Workflow.Upsert(workflow)
	if err != nil {
		app.ErrorLog.Println("Unable to save workflow:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to save workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow saved successfully",
		"workflow_id": id,
	})
}

func (app *Config) triggerWorkflow(c *gin.Context) {

	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	// background update workflow information into database
	go func() {

		nodes, err := json.Marshal(req.Nodes)
		if err != nil {
			app.ErrorLog.Println("Unable to marshal nodes:", err)
			return
		}

		workflow := data.Workflow{
			WorkflowID:   req.WorkflowID,
			WorkflowName: req.WorkflowName,
			RootNodeID:   "start",
			Nodes:        nodes,
		}

		_, err = app.Model.Workflow.Upsert(workflow)
		if err != nil {
			app.ErrorLog.Println("Unable to save workflow:", err)
			return
		}

		app.InfoLog.Println("Workflow saved successfully :", req.WorkflowID)
	}()

	// 2. make sure workflow ID is set
	if req.WorkflowID == "" {
		req.WorkflowID = "robot-routine-" + uuid.New().String()
	}

	app.InfoLog.Println("Received workflow trigger request:", req)

	// TODO: thinking about the request schema different between received from client and interanl useage
	// 3. setting workflow options
	nodesBytes, _ := json.Marshal(req.Nodes)

	// 4. unmarshal req.
	var nodes map[string]pkg.WorkflowNode
	json.Unmarshal(nodesBytes, &nodes)

	var payload = pkg.WorkflowPayload{
		WorkflowID: req.WorkflowID,
		RootNodeID: "start",
		Nodes:      nodes,
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:        payload.WorkflowID,
		TaskQueue: "ROBOT_TASK_QUEUE",
	}

	we, err := app.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RobotWorkflow, payload)
	if err != nil {
		app.ErrorLog.Println("Unable to start workflow:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to start workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow started successfully",
		"workflow_id": we.GetID(),
	})
}

func (app *Config) pauseWorkflow(c *gin.Context) {

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Workflow Id is required"})
		return
	}

	err := app.TemporalClient.SignalWorkflow(context.Background(), workflowID, "", "control-signal", "pause")
	if err != nil {
		app.ErrorLog.Println("Unable to signal workflow:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to signal workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Pause workflow",
		"workflow_id": workflowID,
	})
}

func (app *Config) resumeWorkflow(c *gin.Context) {

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Workflow Id is required"})
		return
	}

	err := app.TemporalClient.SignalWorkflow(context.Background(), workflowID, "", "control-signal", "resume")
	if err != nil {
		app.ErrorLog.Println("Unable to signal workflow:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to signal workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Resume workflow",
		"workflow_id": workflowID,
	})
}

func (app *Config) getWorkflows(c *gin.Context) {
	workflows, err := app.Model.Workflow.Get()
	if err != nil {
		app.ErrorLog.Println("Unable to get workflows:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to get workflows"})
		return
	}

	c.JSON(http.StatusOK, workflows)
}

func (app *Config) getWorkflowById(c *gin.Context) {

	workflow_id := c.Param("id")
	if workflow_id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Workflow Id is required"})
		return
	}

	workflow, err := app.Model.Workflow.GetByID(workflow_id)
	if err != nil {
		app.ErrorLog.Println("Unable to get workflow:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to get workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (app *Config) getWorkflowStatus(c *gin.Context) {

	workflow_id := c.Param("id")
	if workflow_id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Workflow Id is required"})
		return
	}

	// In a real implementation, you would query Temporal for the workflow status.
	resp, err := app.TemporalClient.QueryWorkflow(context.Background(), workflow_id, "", "get_step")
	if err != nil {
		app.ErrorLog.Println("Unable to query workflow:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to query workflow"})
		return
	}

	var currentStep string
	if err := resp.Get(&currentStep); err != nil {
		app.ErrorLog.Println("Unable to get query result:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Unable to get query result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflow_id":  workflow_id,
		"current_step": currentStep,
	})
}
