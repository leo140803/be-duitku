package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leo140803/finance-app-backend/handlers"
	"github.com/leo140803/finance-app-backend/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:3001", "https://fe-duitku-git-main-leonardo-nickholas-andriantos-projects.vercel.app"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	
	r.Use(cors.New(corsConfig))

	api := r.Group("/api")
	{
		//HEALTH CHECK
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "service is healthy",
			})
		})
		// Public routes (no authentication required)
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/logout", handlers.Logout)

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/auth/profile", handlers.GetProfile)

			// Accounts
			protected.GET("/accounts", handlers.GetAccounts)
			protected.POST("/accounts", handlers.CreateAccount)

			// Categories
			protected.GET("/categories", handlers.GetCategories)
			protected.POST("/categories", handlers.CreateCategory)
			protected.PUT("/categories/:id", handlers.UpdateCategory)
			protected.DELETE("/categories/:id", handlers.DeleteCategory)


			// Transactions
			protected.GET("/transactions", handlers.GetTransactions)
			protected.POST("/transactions", handlers.CreateTransaction)
			protected.PUT("/transactions/:id", handlers.UpdateTransaction)
			protected.DELETE("/transactions/:id", handlers.DeleteTransaction)
		}
	}

	return r
}
