package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leo140803/finance-app-backend/config"
	"github.com/leo140803/finance-app-backend/models"
)

func GetCategories(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var categories []models.Category

	err := config.SupaClient.DB.From("categories").Select("*").Eq("user_id", userID.(string)).Execute(context.Background(), &categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set user ID from authentication context
	cat.UserID = userID.(string)
	
	// Clear any existing ID or CreatedAt to let database handle them
	cat.ID = ""
	cat.CreatedAt = ""

	// Use interface{} to handle flexible response format from Supabase
	var result interface{}
	err := config.SupaClient.DB.From("categories").Insert(cat).Execute(context.Background(), &result)
	
	// If insert fails, try to get existing category
	if err != nil {
		// Try to get existing category by name and user_id
		var existingCategories []models.Category
		err = config.SupaClient.DB.From("categories").Select("*").Eq("name", cat.Name).Eq("user_id", cat.UserID).Execute(context.Background(), &existingCategories)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or fetch category: " + err.Error()})
			return
		}
		
		if len(existingCategories) > 0 {
			c.JSON(http.StatusCreated, existingCategories[0])
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category and category not found"})
		}
		return
	}

	// Return success with basic info, frontend will reload the list
	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"name": cat.Name,
		"user_id": cat.UserID,
	})
}

func UpdateCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryID := c.Param("id")

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var updated []models.Category
	err := config.SupaClient.DB.
		From("categories").
		Update(models.Category{
			Name:   input.Name,
			UserID: userID.(string), // supaya filter tetap aman
		}).
		Eq("id", categoryID).
		Eq("user_id", userID.(string)).
		Execute(context.Background(), &updated)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category: " + err.Error()})
		return
	}

	if len(updated) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, updated[0])
}


func DeleteCategory(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get category ID dari URL param
	categoryID := c.Param("id")

	var deleted []models.Category
	err := config.SupaClient.DB.
		From("categories").
		Delete().
		Eq("id", categoryID).
		Eq("user_id", userID.(string)).
		Execute(context.Background(), &deleted)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category: " + err.Error()})
		return
	}

	if len(deleted) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

