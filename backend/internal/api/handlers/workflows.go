package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/chungweeeei/Temporal-robot-project/internal/repository/models"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type SaveWorkflowRequest struct {
	WorkflowID   string                 `json:"workflow_id" binding:"required"`
	WorkflowName string                 `json:"workflow_name" binding:"required"`
	Nodes        map[string]interface{} `json:"nodes"`
}

func (h *Handler) SaveWorkflow(c *gin.Context) {

	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.App.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		h.App.ErrorLog.Println("Unable to marshal nodes:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to marshal nodes"})
		return
	}

	workflow := models.Workflow{
		WorkflowID:   req.WorkflowID,
		WorkflowName: req.WorkflowName,
		RootNodeID:   "start",
		Nodes:        nodes,
	}

	id, err := h.App.Model.Workflow.Upsert(workflow)
	if err != nil {
		h.App.ErrorLog.Println("Unable to save workflow:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to save workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow saved successfully",
		"workflow_id": id,
	})
}

func (h *Handler) TriggerWorkflow(c *gin.Context) {

	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.App.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	// background update workflow information into database
	go func() {

		nodes, err := json.Marshal(req.Nodes)
		if err != nil {
			h.App.ErrorLog.Println("Unable to marshal nodes:", err)
			return
		}

		workflow := models.Workflow{
			WorkflowID:   req.WorkflowID,
			WorkflowName: req.WorkflowName,
			RootNodeID:   "start",
			Nodes:        nodes,
		}

		_, err = h.App.Model.Workflow.Upsert(workflow)
		if err != nil {
			h.App.ErrorLog.Println("Unable to save workflow:", err)
			return
		}

		h.App.InfoLog.Println("Workflow saved successfully :", req.WorkflowID)
	}()

	// 2. make sure workflow ID is set
	if req.WorkflowID == "" {
		req.WorkflowID = "robot-routine-" + uuid.New().String()
	}

	h.App.InfoLog.Println("Received workflow trigger request:", req)

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

	we, err := h.App.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RobotWorkflow, payload)
	if err != nil {
		h.App.ErrorLog.Println("Unable to start workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to start workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow started successfully",
		"workflow_id": we.GetID(),
	})
}

func (h *Handler) PauseWorkflow(c *gin.Context) {

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Workflow Id is required"})
		return
	}

	err := h.App.TemporalClient.SignalWorkflow(context.Background(), workflowID, "", "control-signal", "pause")
	if err != nil {
		h.App.ErrorLog.Println("Unable to signal workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to signal workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Pause workflow",
		"workflow_id": workflowID,
	})
}

func (h *Handler) ResumeWorkflow(c *gin.Context) {

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Workflow Id is required"})
		return
	}

	err := h.App.TemporalClient.SignalWorkflow(context.Background(), workflowID, "", "control-signal", "resume")
	if err != nil {
		h.App.ErrorLog.Println("Unable to signal workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to signal workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Resume workflow",
		"workflow_id": workflowID,
	})
}

func (h *Handler) GetWorkflows(c *gin.Context) {
	workflows, err := h.App.Model.Workflow.Get()
	if err != nil {
		h.App.ErrorLog.Println("Unable to get workflows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get workflows"})
		return
	}

	c.JSON(http.StatusOK, workflows)
}

func (h *Handler) GetWorkflowById(c *gin.Context) {

	workflow_id := c.Param("id")
	if workflow_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Workflow Id is required"})
		return
	}

	workflow, err := h.App.Model.Workflow.GetByID(workflow_id)
	if err != nil {
		h.App.ErrorLog.Println("Unable to get workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *Handler) GetWorkflowStatus(c *gin.Context) {

	workflowId := c.Param("id")
	if workflowId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Workflow Id is required"})
		return
	}

	// In a real implementation, you would query Temporal for the workflow status.
	resp, err := h.App.TemporalClient.QueryWorkflow(context.Background(), workflowId, "", "get_step")
	if err != nil {
		h.App.ErrorLog.Println("Unable to query workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to query workflow"})
		return
	}

	var currentStep string
	if err := resp.Get(&currentStep); err != nil {
		h.App.ErrorLog.Println("Unable to get query result:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get query result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflow_id":  workflowId,
		"current_step": currentStep,
	})
}

func (h *Handler) DeleteWorkflow(c *gin.Context) {

	workflowId := c.Param("id")
	if workflowId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Workflow Id is required"})
		return
	}

	err := h.App.Model.Workflow.Delete(workflowId)
	if err != nil {
		h.App.ErrorLog.Println("Unable to delete workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Workflow deleted successfully",
		"workflow_id": workflowId,
	})
}
