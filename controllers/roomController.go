package controllers

import (
	"my-app/models"
	"net/http"
	"strconv" // For string to int conversion

	"github.com/gin-gonic/gin"
)

// Create a new room (C)
func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.CreateRoom(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Room created successfully"})
}

// Get rooms for a theater (R)
// GetRoomsByTheaterID fetches all rooms for a given theater ID
func GetRoomsByTheaterID(c *gin.Context) {
	theaterIDStr := c.Param("theaterID")
	theaterID, err := strconv.Atoi(theaterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater ID"})
		return
	}

	rooms, err := models.GetRoomsByTheaterID(theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// Update room (U)
func UpdateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.UpdateRoom(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Room updated successfully"})
}

// Delete room (D)
func DeleteRoom(c *gin.Context) {
	roomIDStr := c.Param("id")             // Get roomID as string
	roomID, err := strconv.Atoi(roomIDStr) // Convert to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = models.DeleteRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Room deleted successfully"})
}
