package controllers // Ensure the package is declared

import (
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Function to create a new screen
func CreateScreen(c *gin.Context) {
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := models.CreateScreen(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Screen created successfully", "screen": screen})
}

// Function to get screens by room ID
func GetScreensByRoomID(c *gin.Context) {
	roomIDStr := c.Param("roomID")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	screens, err := models.GetScreensByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"screens": screens})
}

// GetAllScreens fetches all screens
func GetAllScreens(c *gin.Context) {
	screens, err := models.GetAllScreens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"screens": screens})
}

// Function to update an existing screen
func UpdateScreen(c *gin.Context) {
	screenIDStr := c.Param("id")
	screenID, err := strconv.Atoi(screenIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	screen.ScreenID = screenID

	if err := models.UpdateScreen(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Screen updated successfully", "screen": screen})
}

// Function to delete a screen by ID
func DeleteScreen(c *gin.Context) {
	screenIDStr := c.Param("id")
	screenID, err := strconv.Atoi(screenIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	if err := models.DeleteScreen(screenID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Screen deleted successfully"})
}
