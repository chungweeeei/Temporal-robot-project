package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/gin-gonic/gin"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type CreateScheduleRequest struct {
	ScheduleID string `json:"schedule_id" binding:"required"`
	CronExpr   string `json:"cron_expr" binding:"required"` // e.g. "*/5 * * * *"
	Timezone   string `json:"timezone"`                     // 預設 Asia/Taipei
}

type ScheduleInfo struct {
	ScheduleID string `json:"schedule_id"`
	Spec       struct {
		Calendars []ReadableCalendarSpec `json:"calendars"`
	} `json:"spec"`
	State struct {
		Paused           bool   `json:"paused"`
		Note             string `json:"note"`
		RemainingActions int    `json:"remaining_actions"`
	} `json:"state"`
}

type UpdateScheduleRequest struct {
	CronExpr *string `json:"cron_expr"` // e.g. "*/5 * * * *"
}

func (app *Config) createSchedule(c *gin.Context) {

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	// default timezone setting
	timezone := req.Timezone
	if timezone == "" {
		timezone = "Asia/Taipei"
	}

	_, err := time.LoadLocation(timezone)
	if err != nil {
		app.ErrorLog.Println("Invalid timezone:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: "Invalid timezone"})
		return
	}

	// register temporal schedule client
	scheduleClient := app.TemporalClient.ScheduleClient()

	// start creating schedule task
	scheduleHandle, err := scheduleClient.Create(context.Background(), client.ScheduleOptions{
		ID: req.ScheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{req.CronExpr},
			Jitter:          time.Second * 10,
			TimeZoneName:    timezone,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        fmt.Sprintf("%s-workflow", req.ScheduleID),
			Workflow:  workflows.RobotScheduleWorkflow,
			TaskQueue: "ROBOT_SCHEDULE_QUEUE",
		},
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
	})

	if err != nil {
		app.ErrorLog.Println("Unable to create schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to create schedule"})
		return
	}

	app.InfoLog.Printf("Schedule created: %s", scheduleHandle.GetID())

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule created successfully",
		"schedule_id": scheduleHandle.GetID(),
		"cron_expr":   req.CronExpr,
		"timezone":    timezone,
	})
}

func (app *Config) getSchedules(c *gin.Context) {

	scheduleClient := app.TemporalClient.ScheduleClient()

	listView, err := scheduleClient.List(context.Background(), client.ScheduleListOptions{
		PageSize: 1,
	})
	if err != nil {
		app.ErrorLog.Println("Unable to list schedules:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to list schedules"})
		return
	}

	var schedules []ScheduleInfo
	for listView.HasNext() {
		scheduleEntry, err := listView.Next()
		if err != nil {
			app.ErrorLog.Println("Error iterating schedules:", err)
			break
		}
		schedules = append(schedules, ScheduleInfo{
			ScheduleID: scheduleEntry.ID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
	})
}

func (app *Config) getScheduleById(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	scheduleInfo, err := scheduleHandle.Describe(context.Background())
	if err != nil {
		app.ErrorLog.Println("Unable to describe schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to describe schedule"})
		return
	}

	var readableCalendars = []ReadableCalendarSpec{}
	for _, spec := range scheduleInfo.Schedule.Spec.Calendars {
		readableCalendars = append(readableCalendars, ReadableCalendarSpec{
			Second:  rangeToString(spec.Second),
			Minute:  rangeToString(spec.Minute),
			Hour:    rangeToString(spec.Hour),
			Month:   rangeToString(spec.Month),
			Year:    rangeToString(spec.Year),
			Comment: spec.Comment,
		})
	}

	var response ScheduleInfo
	response.ScheduleID = scheduleID
	response.Spec.Calendars = readableCalendars
	response.State.Paused = scheduleInfo.Schedule.State.Paused
	response.State.Note = scheduleInfo.Schedule.State.Note
	response.State.RemainingActions = scheduleInfo.Schedule.State.RemainingActions

	c.JSON(http.StatusOK, response)
}

func (app *Config) triggerSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Trigger(context.Background(), client.ScheduleTriggerOptions{
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,
	})
	if err != nil {
		app.ErrorLog.Println("Unable to trigger schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to trigger schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule triggered successfully",
		"schedule_id": scheduleID,
	})
}

func (app *Config) pauseSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Pause(context.Background(), client.SchedulePauseOptions{
		Note: "The Schedule has been paused.",
	})
	if err != nil {
		app.ErrorLog.Println("Unable to pause schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to pause schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule paused successfully",
		"schedule_id": scheduleID,
	})
}

func (app *Config) resumeSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Unpause(context.Background(), client.ScheduleUnpauseOptions{
		Note: "The Schedule has been resumed.",
	})
	if err != nil {
		app.ErrorLog.Println("Unable to resume schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to resume schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule resumed successfully",
		"schedule_id": scheduleID,
	})
}

func (app *Config) deleteSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()

	scheduleHandle := scheduleClient.GetHandle(c, scheduleID)
	err := scheduleHandle.Delete(context.Background())
	if err != nil {
		app.ErrorLog.Println("Unable to delete schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule deleted successfully",
		"schedule_id": scheduleID,
	})
}

func (app *Config) updateSchedule(c *gin.Context) {

	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Schedule Id is required"})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.ErrorLog.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Message: fmt.Sprintf("Invalid payload: %v", err)})
		return
	}

	scheduleClient := app.TemporalClient.ScheduleClient()
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
		app.ErrorLog.Println("Unable to update schedule:", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{Message: "Unable to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Schedule updated successfully",
		"schedule_id": scheduleID,
	})
}
