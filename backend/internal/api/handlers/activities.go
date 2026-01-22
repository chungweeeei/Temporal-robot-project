package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetActivities(c *gin.Context) {

	activities, err := h.App.Model.Activity.Get()
	if err != nil {
		h.App.ErrorLog.Println("Unable to get activities:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}
