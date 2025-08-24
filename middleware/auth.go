package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leo140803/finance-app-backend/config"
	"github.com/leo140803/finance-app-backend/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify token with Supabase Auth
		authUser, err := config.SupaClient.Auth.User(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Get user from our users table using email from auth user
		var users []models.User
		err = config.SupaClient.DB.From("users").Select("*").Eq("email", authUser.Email).Execute(context.Background(), &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
			c.Abort()
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found in database"})
			c.Abort()
			return
		}

		// Set user ID from our database (not from Supabase Auth)
		c.Set("user_id", users[0].ID)
		c.Next()
	}
}
