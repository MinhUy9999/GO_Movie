package controllers

import (
	"my-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book tickets (C)
func BookTickets(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Kiểm tra kiểu dữ liệu của userID
	var userID int
	switch v := userIDInterface.(type) {
	case int:
		userID = v
	case uint:
		userID = int(v)
	case float64:
		userID = int(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	var request struct {
		ScheduleID int   `json:"schedule_id"`
		Seats      []int `json:"seats"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	bookingID, err := models.BookSeats(userID, request.ScheduleID, request.Seats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Booking successful", "booking_id": bookingID})
}
