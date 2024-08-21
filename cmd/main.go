package main

import (
	"github-service/internal/config"
	"github-service/internal/database"
	"github-service/internal/routes"
	"github-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database

	db := database.Connect(config.GetDatabaseConfig("development"))
	//gin.SetMode(gin.ReleaseMode)
	// Create Gin router
	router := gin.Default()

	// Register handlers
	routes.APPRoutes(router, db)

	// Start the background task for continuous data fetching
	go service.StartDataFetchingTask(db, "chromium")
	// Run the server
	router.Run(":8080")
}
