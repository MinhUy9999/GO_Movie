package controllers

import (
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookTickets books seats
func BookTickets(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

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

// GetBookingsByUserID retrieves all bookings for a given user
func GetBookingsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bookings, err := models.GetBookingDetailsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success", "bookings": bookings})
}

// GetAllBookings retrieves all bookings
func GetAllBookings(c *gin.Context) {
	bookings, err := models.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success", "bookings": bookings})
}

// DeleteBooking deletes a booking by ID
func DeleteBooking(c *gin.Context) {
	bookingIDParam := c.Param("booking_id")
	bookingID, err := strconv.ParseInt(bookingIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = models.DeleteBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Booking deleted successfully"})
}
