package main

import (
	"context"
	"github-service/internal/config"
	"github-service/internal/database"
	"github-service/internal/database/repository"
	"github-service/internal/service"
	"github-service/internal/web/handlers"
	"github-service/internal/web/routes"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation when done

	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the database connection
	db, err := database.Connect(config.GetDatabaseConfig(cfg, "dev"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a new Gin router instance
	router := gin.Default()

	// Initialize repository instances
	commitRepo, err := repository.NewCommitRepository(db)
	if err != nil {
		log.Fatalf("Failed to create commit repository: %v", err)
	}
	repositoryRepo, err := repository.NewRepository(db)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize service instances
	commitService := service.NewCommitService(commitRepo, cfg)
	repositoryService := service.NewRepositoryService(repositoryRepo, cfg)

	// Initialize handlers
	commitHandler := handlers.NewCommitHandler(commitService, repositoryService)
	repositoryHandler := handlers.NewRepositoryHandler(repositoryService)

	// Initialize the commit monitor service
	commitMonitor := service.NewCommitMonitor(commitService, repositoryService)

	// Initialize API routes
	routes.SetupAPIRoutes(router, commitHandler, repositoryHandler)

	// Retrieve configuration from environment variables
	repoOwner := cfg.SEED_REPO_OWNER
	repoName := cfg.SEED_REPO_NAME
	beginFetchDateStr := cfg.BEGIN_FETCH_DATE

	// Parse the beginFetchDate from the environment variable
	beginFetchDate, err := time.Parse(time.RFC3339, beginFetchDateStr)
	if err != nil {
		log.Fatalf("Invalid BEGIN_FETCH_DATE format: %v", err)
	}

	// Seed the database with initial data starting from the defined date
	commitMonitor.SeedDB(repoOwner, repoName, beginFetchDate)

	// Start the background task that periodically fetches and stores repository data
	go commitMonitor.StartDataFetchingTask(ctx, repoOwner, repoName)

	// Start the Gin HTTP server and listen on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
