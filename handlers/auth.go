package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leo140803/finance-app-backend/config"
	"github.com/leo140803/finance-app-backend/models"
	"github.com/lengzuo/supa/dto"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Register user with Supabase Auth
	signUpReq := dto.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	
	authResponse, err := config.SupaClient.Auth.SignUp(context.Background(), signUpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	// Create user record in our users table
	user := models.User{
		Email: req.Email,
	}

	// Try to insert user record, but don't fail if it already exists
	var result []models.User
	err = config.SupaClient.DB.From("users").Insert(user).Execute(context.Background(), &result)
	
	// If insert fails, try to get existing user
	if err != nil {
		// Try to get existing user by email
		var existingUsers []models.User
		err = config.SupaClient.DB.From("users").Select("*").Eq("email", req.Email).Execute(context.Background(), &existingUsers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or fetch user record: " + err.Error()})
			return
		}
		
		if len(existingUsers) > 0 {
			result = existingUsers
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record and user not found"})
			return
		}
	}

	response := models.AuthResponse{
		User:         result[0],
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Login user with Supabase Auth
	signInReq := dto.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	
	authResponse, err := config.SupaClient.Auth.SignInWithPassword(context.Background(), signInReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: " + err.Error()})
		return
	}

	// Get user from our users table
	var users []models.User
	err = config.SupaClient.DB.From("users").Select("*").Eq("email", req.Email).Execute(context.Background(), &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := models.AuthResponse{
		User:         users[0],
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var users []models.User
	err := config.SupaClient.DB.From("users").Select("*").Eq("id", userID.(string)).Execute(context.Background(), &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile: " + err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, users[0])
}

func Logout(c *gin.Context) {
	// Get refresh token from request
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	// Logout user with Supabase Auth
	err := config.SupaClient.Auth.SignOut(context.Background(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
