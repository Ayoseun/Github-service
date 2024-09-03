package main

import (
	"context"
	"github-service/internal/config"
	"github-service/internal/database"
	"github-service/internal/database/repository"
	"github-service/internal/service"
	"github-service/internal/web/handlers"
	"github-service/internal/web/routes"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation when done
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Initialize the database connection
	// Connect to the database using configuration settings for the "dev" environment
	db, err := database.Connect(config.GetDatabaseConfig(cfg, "dev"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Create a new Gin router instance
	router := gin.Default()

	// Initialize repository
	commitRepo, err := repository.NewCommitRepository(db)
	if err != nil {

		log.Fatalf("Failed to create repository: database connection is nil%v", err)
	}
	repositoryRepository, err := repository.NewRepository(db)
	if err != nil {
		log.Fatalf("ailed to create repository: database connection is nil: %v", err)
	}
	// Initialize service
	commitService := service.NewCommitService(commitRepo, cfg)
	repositoryService := service.NewRepositoryService(repositoryRepository, cfg)
	commitHandler := handlers.NewCommitHandler(commitService, repositoryService)
	repositoryHandler := handlers.NewRepositoryHandler(repositoryService)
	commitMonitor := service.NewCommitMonitor(commitService, repositoryService)
	// Initialize routes
	routes.SetupAPIRoutes(router, commitHandler, repositoryHandler)
	// Define the initial fetch date for seeding the database
	beginFetchDate := time.Date(2022, 12, 9, 0, 0, 0, 0, time.UTC) // Example date: December 9, 2022
	// Seed the database with initial data starting from the defined date
	commitMonitor.SeedDB("golang", "go", beginFetchDate)

	// Start the background task that periodically fetches and stores repository data

	go commitMonitor.StartDataFetchingTask(ctx, cfg, "golang", "go")

	// Start the Gin HTTP server and listen on port 8080
	router.Run(":8080")

}
