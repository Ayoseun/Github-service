package main

import (
	"github-service/internal/config"           // Import the package for configuration management
	"github-service/internal/database/storage" // Import the package for database storage functions
	"github-service/internal/domain/service"   // Import the package for domain-specific services
	"github-service/internal/web/routes"       // Import the package for defining web routes
	"time"                                     // Import the time package for handling dates and times

	"github.com/gin-gonic/gin" // Import the Gin framework for routing and HTTP server functionalities
)

func main() {
	// Initialize the database connection
	// Connect to the database using configuration settings for the "dev" environment
	db := storage.Connect(config.GetDatabaseConfig("dev"))

	// Create a new Gin router instance
	router := gin.Default()

	// Register application routes with the Gin router and the database connection
	routes.APPRoutes(router, db)

	// Define the initial fetch date for seeding the database
	beginFetchDate := time.Date(2022, 12, 9, 0, 0, 0, 0, time.UTC) // Example date: December 9, 2022

	// Seed the database with initial data starting from the defined date
	service.SeedDB(db, "chromium", beginFetchDate)
	// Start the background task that periodically fetches and stores repository data
	go service.StartDataFetchingTask(db, "chromium")
	// Start the Gin HTTP server and listen on port 8080
	router.Run(":8080")

}
