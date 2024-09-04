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
	"os"
	"os/signal"
	"syscall"
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

	// Initialize the commit monitor service
	commitMonitor := service.NewCommitMonitor(ctx, commitService, repositoryService)

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
	if _, err := commitMonitor.AddRepository(repoOwner, repoName, beginFetchDate); err != nil {

		log.Printf("Failed to add initial repository: %v", err)

	}
	// Start the background task that periodically fetches and stores repository data
	go commitMonitor.StartDataFetchingTask(ctx, repoOwner, repoName)

	// Handle graceful shutdown
	go func() {
		// Create a channel to listen for OS signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// Block until a signal is received
		<-sigs
		log.Println("Shutdown signal received, stopping the commit monitor...")
		cancel() // Cancel the context to stop the background task
	}()

	// Initialize the Gin router and API routes
	router := gin.Default()
	commitHandler := handlers.NewCommitHandler(commitService, repositoryService)
	repositoryHandler := handlers.NewRepositoryHandler(repositoryService, commitMonitor)
	routes.SetupAPIRoutes(router, commitHandler, repositoryHandler)

	// Start the Gin HTTP server and listen on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
