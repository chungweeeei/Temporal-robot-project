package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/internal/workflow"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gin-gonic/gin"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type CreateScheduleRequest struct {
	ScheduleID string `json:"schedule_id" binding:"required"`
	WorkflowID string `json:"workflow_id" binding:"required"`
	CronExpr   string `json:"cron_expr" binding:"required"` // e.g. "*/5 * * * *"
	Timezone   string `json:"timezone"`                     // 預設 Asia/Taipei
}

type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Step  int `json:"step"`
}

type CalendarSpec struct {
	Second     []Range `json:"second,omitempty"`
	Minute     []Range `json:"minute,omitempty"`
	Hour       []Range `json:"hour,omitempty"`
	DayOfMonth []Range `json:"day_of_month,omitempty"`
	Month      []Range `json:"month,omitempty"`
	Year       []Range `json:"year,omitempty"`
	DayOfWeek  []Range `json:"day_of_week,omitempty"`
	Comment    string  `json:"comment,omitempty"`
}

type ScheduleInfo struct {
	ScheduleID string `json:"schedule_id"`
	Spec       struct {
		Calendars       []CalendarSpec `json:"calendars"`
		CronExpressions []string       `json:"cron_expressions"`
	} `json:"spec"`
	Paused      bool   `json:"paused"`
	RecentRun   string `json:"recent_run"`
	UpcomingRun string `json:"upcoming_run"`
}

type UpdateScheduleRequest struct {
	CronExpr *string `json:"cron_expr"` // e.g. "*/5 * * * *"
}

func transformToRanges(ranges []client.ScheduleRange) []Range {

	var result []Range
	if len(ranges) == 0 {
		return []Range{}
	}

	for _, r := range ranges {
		end := r.End
		if end == 0 {
			end = r.Start
		}
		step := r.Step
		if step == 0 {
			step = 1
		}
		result = append(result, Range{
			Start: r.Start,
			End:   end,
			Step:  step,
		})
	}
	return result
}

func (h *Handler) CreateSchedule(c *gin.Context) {

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.App.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	// default timezone setting
	timezone := req.Timezone
	if timezone == "" {
		timezone = "Asia/Taipei"
	}

	_, err := time.LoadLocation(timezone)
	if err != nil {
		h.App.ErrorLog.Println("Invalid timezone:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "Invalid timezone"})
		return
	}

	// check workflow id existence
	record, err := h.App.Model.Workflow.GetByID(req.WorkflowID)
	if err != nil {
		h.App.ErrorLog.Println("Unable to get workflow:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get workflow"})
		return
	}

	var nodes map[string]pkg.WorkflowNode
	if err := json.Unmarshal([]byte(record.Nodes), &nodes); err != nil {
		h.App.ErrorLog.Println("Unable to unmarshal nodes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to process workflow data"})
		return
	}

	// register temporal schedule client
	scheduleClient := h.App.TemporalClient.ScheduleClient()

	// start creating schedule task
	scheduleHandle, err := scheduleClient.Create(context.Background(), client.ScheduleOptions{
		ID: req.ScheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{req.CronExpr},
			Jitter:          time.Second * 10,
			TimeZoneName:    timezone,
		},
		Action: &client.ScheduleWorkflowAction{
			// ID: 當schedule啟動workflow時產生的workflowID, ex: test-workflow2-schedule-001-2026-01-20T02:06:33Z
			ID: req.ScheduleID,
			// 註冊在 Temporal server 的 Workflow 名稱
			Workflow: workflow.RobotWorkflow,
			// 如果 TaskQueue 也是存在 DB，可以用 record.TaskQueue，否則這裡是寫死的
			TaskQueue: "ROBOT_TASK_QUEUE",
			Args: []interface{}{pkg.WorkflowPayload{
				WorkflowID: record.WorkflowID,
				RootNodeID: record.RootNodeID,
				Nodes:      nodes,
			}},
		},
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
	})

	if err != nil {
		h.App.ErrorLog.Println("Unable to create schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to create schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule created successfully",
		"schedule_id": scheduleHandle.GetID(),
		"cron_expr":   req.CronExpr,
		"timezone":    timezone,
	})
}

func (h *Handler) GetSchedules(c *gin.Context) {

	scheduleClient := h.App.TemporalClient.ScheduleClient()

	listView, err := scheduleClient.List(context.Background(), client.ScheduleListOptions{
		PageSize: 1,
	})
	if err != nil {
		h.App.ErrorLog.Println("Unable to list schedules:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to list schedules"})
		return
	}

	schedules := []ScheduleInfo{}
	for listView.HasNext() {
		scheduleEntry, err := listView.Next()
		if err != nil {
			h.App.ErrorLog.Println("Error iterating schedules:", err)
			break
		}

		var calendars = []CalendarSpec{}
		if scheduleEntry.Spec != nil {
			for _, spec := range scheduleEntry.Spec.Calendars {
				calendars = append(calendars, CalendarSpec{
					Second:     transformToRanges(spec.Second),
					Minute:     transformToRanges(spec.Minute),
					Hour:       transformToRanges(spec.Hour),
					DayOfMonth: transformToRanges(spec.DayOfMonth),
					Month:      transformToRanges(spec.Month),
					Year:       transformToRanges(spec.Year),
					DayOfWeek:  transformToRanges(spec.DayOfWeek),
					Comment:    spec.Comment,
				})
			}
		}

		info := ScheduleInfo{
			ScheduleID: scheduleEntry.ID,
			Paused:     scheduleEntry.Paused,
		}
		if scheduleEntry.Spec != nil {
			info.Spec.Calendars = calendars
			info.Spec.CronExpressions = scheduleEntry.Spec.CronExpressions
		}

		if len(scheduleEntry.RecentActions) > 0 {
			info.RecentRun = scheduleEntry.RecentActions[len(scheduleEntry.RecentActions)-1].ScheduleTime.Format(time.RFC3339)
		}

		if len(scheduleEntry.NextActionTimes) > 0 {
			info.UpcomingRun = scheduleEntry.NextActionTimes[0].Format(time.RFC3339)
		}

		schedules = append(schedules, info)
	}

	c.JSON(http.StatusOK, schedules)
}

func (h *Handler) GetScheduleById(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Schedule Id is required"})
		return
	}

	scheduleClient := h.App.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	scheduleInfo, err := scheduleHandle.Describe(context.Background())
	if err != nil {
		h.App.ErrorLog.Println("Unable to describe schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to describe schedule"})
		return
	}

	var calendars = []CalendarSpec{}
	for _, spec := range scheduleInfo.Schedule.Spec.Calendars {
		calendars = append(calendars, CalendarSpec{
			Second:     transformToRanges(spec.Second),
			Minute:     transformToRanges(spec.Minute),
			Hour:       transformToRanges(spec.Hour),
			DayOfMonth: transformToRanges(spec.DayOfMonth),
			Month:      transformToRanges(spec.Month),
			Year:       transformToRanges(spec.Year),
			DayOfWeek:  transformToRanges(spec.DayOfWeek),
			Comment:    spec.Comment,
		})
	}

	var response ScheduleInfo
	response.ScheduleID = scheduleID
	response.Spec.Calendars = calendars
	response.Spec.CronExpressions = scheduleInfo.Schedule.Spec.CronExpressions
	response.Paused = scheduleInfo.Schedule.State.Paused

	if len(scheduleInfo.Info.RecentActions) > 0 {
		lastAction := scheduleInfo.Info.RecentActions[len(scheduleInfo.Info.RecentActions)-1]
		response.RecentRun = lastAction.ScheduleTime.Format(time.RFC3339)
	}

	if len(scheduleInfo.Info.NextActionTimes) > 0 {
		nextAction := scheduleInfo.Info.NextActionTimes[0]
		response.UpcomingRun = nextAction.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) PauseSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Schedule Id is required"})
		return
	}

	scheduleClient := h.App.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Pause(context.Background(), client.SchedulePauseOptions{
		Note: "The Schedule has been paused.",
	})
	if err != nil {
		h.App.ErrorLog.Println("Unable to pause schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to pause schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule paused successfully",
		"schedule_id": scheduleID,
	})
}

func (h *Handler) ResumeSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Schedule Id is required"})
		return
	}

	scheduleClient := h.App.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Unpause(context.Background(), client.ScheduleUnpauseOptions{
		Note: "The Schedule has been resumed.",
	})
	if err != nil {
		h.App.ErrorLog.Println("Unable to resume schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to resume schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule resumed successfully",
		"schedule_id": scheduleID,
	})
}

func (h *Handler) DeleteSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Schedule Id is required"})
		return
	}

	scheduleClient := h.App.TemporalClient.ScheduleClient()
	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Delete(context.Background())
	if err != nil {
		h.App.ErrorLog.Println("Unable to delete schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule deleted successfully",
		"schedule_id": scheduleID,
	})
}

func (h *Handler) UpdateSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Schedule Id is required"})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.App.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	scheduleClient := h.App.TemporalClient.ScheduleClient()
	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)

	// define the update function
	updateSchedule := func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {

		schedule := input.Description.Schedule

		if req.CronExpr != nil {
			schedule.Spec.CronExpressions = []string{*req.CronExpr}
			schedule.Spec.Calendars = nil
			schedule.Spec.Intervals = nil
		}

		return &client.ScheduleUpdate{
			Schedule: &schedule,
		}, nil
	}

	err := scheduleHandle.Update(context.Background(), client.ScheduleUpdateOptions{
		DoUpdate: updateSchedule,
	})
	if err != nil {
		h.App.ErrorLog.Println("Unable to update schedule:", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unable to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule updated successfully",
		"schedule_id": scheduleID,
	})
}
