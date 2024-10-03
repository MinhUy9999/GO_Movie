package controllers

import (
	"my-app/models"
	"my-app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set default role to "user"
	user.Role = "user"

	err := models.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Registration successful"})
}

func Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, isAuthenticated, err := models.AuthenticateUser(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if isAuthenticated {
		token, err := utils.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

// thay đổi
// UserProfile - accessible by any authenticated user (user or admin)
func UserProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	c.JSON(http.StatusOK, gin.H{"status": "Welcome!", "user_id": userID, "role": role})
}

// AdminDashboard - accessible only by admins
func AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Admin Dashboard Access"})
}

//den day
