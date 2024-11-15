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

// GetAllSeats retrieves all seats
func GetAllSeats(c *gin.Context) {
	seats, err := models.GetAllSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"seats": seats})
}

// UpdateSeat updates an existing seat
// UpdateSeatHandler - Handler for updating a seat
func UpdateSeatHandler(c *gin.Context) {
	seatIDStr := c.Param("id")
	seatID, err := strconv.Atoi(seatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seat ID"})
		return
	}

	var seat models.Seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	seat.SeatID = seatID // Ensure the seat ID from the URL is assigned

	if err := models.UpdateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat updated successfully"})
}

// DeleteSeat deletes a seat by its ID
func DeleteSeat(c *gin.Context) {
	seatIDStr := c.Param("id")
	seatID, err := strconv.Atoi(seatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seat ID"})
		return
	}

	if err := models.DeleteSeat(seatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat deleted successfully"})
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
