package main

import (
	"github-service/internal/config"
	"github-service/internal/database"
	"github-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db := database.Connect(config.GetDatabaseConfig())
	//gin.SetMode(gin.ReleaseMode)
	// Create Gin router
	router := gin.Default()

	// Register handlers
	routes.APPRoutes(router, db)
	// Run the server
	router.Run(":8080")
}
