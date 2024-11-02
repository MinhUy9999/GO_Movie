package controllers

import (
	"fmt"
	"log"
	"net/http"

	"my-app/config"
	"my-app/models"
	"my-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const MAX_LOGIN_ATTEMPTS = 5

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
	var loginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Attempting login for user: %s", loginInput.Email)

	// Check login attempts
	attempts, err := checkLoginAttempts(c, loginInput.Email)
	if err != nil {
		log.Printf("Error checking login attempts: %v", err)
		attempts = 0
	}

	log.Printf("Login attempts for %s: %d", loginInput.Email, attempts)

	if attempts >= MAX_LOGIN_ATTEMPTS {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts. Please try again later."})
		return
	}

	// Authenticate user
	user, authenticated, err := models.AuthenticateUser(loginInput.Email, loginInput.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !authenticated {
		log.Printf("Authentication failed for user %s", loginInput.Email)
		if err := incrementLoginAttempts(c, loginInput.Email); err != nil {
			log.Printf("Error incrementing login attempts: %v", err)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Reset login attempts after successful login
	if err := resetLoginAttempts(c, loginInput.Email); err != nil {
		log.Printf("Error resetting login attempts: %v", err)
	}

	// Generate token
	token, err := utils.GenerateToken(uint(user.ID), user.Role, user.Name)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "name": user.Name, "role": user.Role, "userID": user.ID})
}

func checkLoginAttempts(c *gin.Context, email string) (int, error) {
	redisClient := config.RedisClient
	if redisClient == nil {
		log.Println("Redis client is nil")
		return 0, nil
	}

	key := fmt.Sprintf("login_attempts:%s", email)
	val, err := redisClient.Get(c.Request.Context(), key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("error getting login attempts: %v", err)
	}
	return val, nil
}

func incrementLoginAttempts(c *gin.Context, email string) error {
	redisClient := config.RedisClient
	if redisClient == nil {
		log.Println("Redis client is nil, skipping increment")
		return nil
	}

	key := fmt.Sprintf("login_attempts:%s", email)
	_, err := redisClient.Incr(c.Request.Context(), key).Result()
	if err != nil {
		return fmt.Errorf("error incrementing login attempts: %v", err)
	}
	return nil
}

func resetLoginAttempts(c *gin.Context, email string) error {
	redisClient := config.RedisClient
	if redisClient == nil {
		log.Println("Redis client is nil, skipping reset")
		return nil
	}

	key := fmt.Sprintf("login_attempts:%s", email)
	_, err := redisClient.Del(c.Request.Context(), key).Result()
	if err != nil {
		return fmt.Errorf("error resetting login attempts: %v", err)
	}
	return nil
}

// UserProfile - accessible by any authenticated user (user or admin)
func UserProfile(c *gin.Context) {

	// Check token in Redis
	redisClient := config.RedisClient
	token := c.GetHeader("Authorization")
	_, err := redisClient.Get(c.Request.Context(), token).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Welcome!"})
}

// AdminDashboard - accessible only by admins
func AdminDashboard(c *gin.Context) {
	// Check token in Redis
	redisClient := config.RedisClient
	token := c.GetHeader("Authorization")
	_, err := redisClient.Get(c.Request.Context(), token).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Admin Dashboard Access"})
}
