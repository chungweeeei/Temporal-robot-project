package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *Config) routes() http.Handler {

	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/workflows", app.saveWorkflow)
		apiV1.GET("/workflows", app.getWorkflows)
		apiV1.GET("/workflows/:id", app.getWorkflowById)
		apiV1.GET("/workflows/:id/status", app.getWorkflowStatus)
		apiV1.POST("/workflows/:id/trigger", app.triggerWorkflow)
	}

	return e
}
