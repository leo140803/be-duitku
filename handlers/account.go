package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leo140803/finance-app-backend/config"
	"github.com/leo140803/finance-app-backend/models"
)

func GetAccounts(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var accounts []models.Account
	
	err := config.SupaClient.DB.From("accounts").Select("*").Eq("user_id", userID.(string)).Execute(context.Background(), &accounts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func CreateAccount(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var acc models.Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set user ID from authentication context
	acc.UserID = userID.(string)
	
	// Clear any existing ID or CreatedAt to let database handle them
	acc.ID = ""
	acc.CreatedAt = ""

	// Use interface{} to handle flexible response format from Supabase
	var result interface{}
	err := config.SupaClient.DB.From("accounts").Insert(acc).Execute(context.Background(), &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success with basic info, frontend will reload the list
	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"name": acc.Name,
		"initial_balance": acc.InitialBalance,
		"user_id": acc.UserID,
	})
}
