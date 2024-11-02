package controllers

import (
	"net/http"
	"strconv"

	"my-app/models"

	"github.com/gin-gonic/gin"
)

// CreateScheduleHandler - Handler for creating a new schedule
func CreateScheduleHandler(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateSchedule(&schedule); err != nil {
		// Log the actual error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// GetSchedulesHandler - Handler for retrieving all schedules
func GetSchedulesHandler(c *gin.Context) {
	schedules, err := models.GetSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByIDHandler - Get schedule by ID
func GetScheduleByIDHandler(c *gin.Context) {
	idStr := c.Param("id") // Get the id as a string

	// Convert the string id to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	schedule, err := models.GetScheduleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateScheduleHandler - Handler for updating a schedule
func UpdateScheduleHandler(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// DeleteScheduleHandler - Delete schedule by ID
func DeleteScheduleHandler(c *gin.Context) {
	idStr := c.Param("id") // Get the id as a string

	// Convert the string id to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = models.DeleteSchedule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // No content response for successful deletion
}

// GetSchedulesByScreenIDHandler - Handler for retrieving schedules by screenID
func GetSchedulesByScreenIDHandler(c *gin.Context) {
	screenIDStr := c.Param("screenID")
	screenID, err := strconv.Atoi(screenIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	schedules, err := models.GetSchedulesByScreenID(screenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}
