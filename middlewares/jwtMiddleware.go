package middlewares

import (
	"net/http"
	"strings"

	"my-app/config"
	"my-app/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates tokens and checks roles
func JWTAuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check token in Redis
		redisClient := config.RedisClient
		if redisClient != nil {
			_, err = redisClient.Get(c.Request.Context(), tokenStr).Result()
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}
		}

		// If requiredRole is set, check that the user has the correct role
		if requiredRole != "" && claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Set user ID and role to the request context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func UseRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redisClient", config.RedisClient)
		c.Next()
	}
}
