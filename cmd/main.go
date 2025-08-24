package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/leo140803/finance-app-backend/config"
	"github.com/leo140803/finance-app-backend/routes"
)

func main() {
	godotenv.Load()

	config.InitDB()

	r := routes.SetupRouter()
	port := os.Getenv("PORT")

	log.Println("Server running on port " + port)
	r.Run(":" + port)
}
