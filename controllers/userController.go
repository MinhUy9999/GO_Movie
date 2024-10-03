package controllers

import (
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Update user details by ID (U)// Update user details by ID (U)
func UpdateUserByID(c *gin.Context) {
	idParam := c.Param("id")             // Get the ID from the URL
	userID, err := strconv.Atoi(idParam) // Convert it to an integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user.ID = userID // Assign the ID from the URL to the user
	err = models.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User updated successfully"})
}

// Delete user by ID (D)
func DeleteUserByID(c *gin.Context) {
	idParam := c.Param("id")             // Get the ID from the URL
	userID, err := strconv.Atoi(idParam) // Convert it to an integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = models.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "User deleted successfully"})
}

// Get all users (Admin only)
func GetAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
