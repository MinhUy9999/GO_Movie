package controllers // Ensure the package is declared

import (
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateSeat handles the creation of a new seat
func CreateSeat(c *gin.Context) {
	var seat models.Seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := models.CreateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seat created successfully", "seat": seat})
}

// Example function to get seats (R)
func GetSeatsByScreenID(c *gin.Context) {
	screenIDStr := c.Param("screenID")
	screenID, err := strconv.Atoi(screenIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	seats, err := models.GetSeatsByScreenID(screenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seats": seats})
}
