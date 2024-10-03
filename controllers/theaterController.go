package controllers

import (
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Create a new theater (C)
func CreateTheater(c *gin.Context) {
	var theater models.Theater
	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.CreateTheater(theater)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Theater created successfully"})
}

// Get theaters (R)
func GetTheaters(c *gin.Context) {
	theaters, err := models.GetAllTheaters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"theaters": theaters})
}

func UpdateTheater(c *gin.Context) {
	// Get the theater ID from the URL
	idParam := c.Param("id")
	theaterID, err := strconv.Atoi(idParam) // Convert to integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater ID"})
		return
	}

	// Bind the updated data to the Theater struct
	var theater models.Theater
	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set the theater ID from the URL
	theater.TheaterID = theaterID

	// Call the model function to update the theater
	err = models.UpdateTheater(theater)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Theater updated successfully"})
}

// Delete theater (D)
func DeleteTheater(c *gin.Context) {
	theaterIDStr := c.Param("id")                // Get theaterID as string
	theaterID, err := strconv.Atoi(theaterIDStr) // Convert theaterID to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater ID"})
		return
	}

	err = models.DeleteTheater(theaterID) // Pass the converted int to the model
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Theater deleted successfully"})
}
