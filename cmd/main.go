package main

import (
	"github-service/internal/config"
	"github-service/internal/database/storage"
	"github-service/internal/domain/service"
	"github-service/internal/web/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database

	db := storage.Connect(config.GetDatabaseConfig("development"))
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
