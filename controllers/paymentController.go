package controllers

import (
	"my-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Process payment (C)
func ProcessPayment(c *gin.Context) {
	userID := c.GetInt("user_id") // Assuming you have user_id from a middleware or JWT context
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment details"})
		return
	}

	// Call the ProcessPayment function from the models package
	err := models.ProcessPayment(userID, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Payment successful"})
}
