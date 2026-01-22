package api

import (
	"time"

	api "github.com/chungweeeei/Temporal-robot-project/internal/api/handlers"
	config "github.com/chungweeeei/Temporal-robot-project/internal/config/api"
	"github.com/chungweeeei/Temporal-robot-project/internal/repository/dao"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *config.AppConfig) *gin.Engine {

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	h := api.NewHandler(app, dao.NewWorkflowDAO(app.DB))

	apiV1 := router.Group("/api/v1")
	{
		// Activities
		apiV1.GET("/activities", h.GetActivities)

		// Workflows for manually trigger
		apiV1.POST("/workflows", h.SaveWorkflow)
		apiV1.GET("/workflows", h.GetWorkflows)
		apiV1.GET("/workflows/records", h.GetWorkflowRecords)
		apiV1.GET("/workflows/:id", h.GetWorkflowById)
		apiV1.GET("/workflows/:id/status", h.GetWorkflowStatus)
		apiV1.POST("/workflows/:id/trigger", h.TriggerWorkflow)
		apiV1.POST("/workflows/:id/pause", h.PauseWorkflow)
		apiV1.POST("/workflows/:id/resume", h.ResumeWorkflow)
		apiV1.DELETE("/workflows/:id", h.DeleteWorkflow)

		// Schedules for scheduled trigger
		apiV1.POST("/schedules", h.CreateSchedule)
		apiV1.GET("/schedules", h.GetSchedules)
		apiV1.GET("/schedules/:id", h.GetScheduleById)
		apiV1.POST("/schedules/:id/pause", h.PauseSchedule)
		apiV1.POST("/schedules/:id/resume", h.ResumeSchedule)
		apiV1.DELETE("/schedules/:id", h.DeleteSchedule)
		apiV1.PUT("/schedules/:id", h.UpdateSchedule)
	}

	return router
}
