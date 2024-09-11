package main

import (
	"context"
	"fmt"
	"log"

	"github-service/config"
	"github-service/internal/adapters"
	"github-service/internal/core/domain"
	"github-service/internal/core/service"
	"github-service/internal/web/handlers"
	"github-service/internal/web/routes"
	"github-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a cancellable context to manage lifecycle and cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context cancellation on exit

	// Load the application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println(cfg) // Print configuration for debugging purposes

	// Initialize the logger
	logger.InitLogger()

	// Set up storage components (Postgres and BadgerDB)
	commitRepo, repositoryRepo, badgerService, err := adapters.SetupStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to setup storage: %v", err)
	}

	// Initialize domain data for default repository
	defaultRepoData := domain.RepoData{
		RepoName: cfg.DEFAULT_REPO,
		Owner:    cfg.DEFAULT_OWNER,
	}

	// Set up core services with the initialized repositories
	commitService, repositoryService, monitorService := service.SetupService(ctx, cfg, defaultRepoData, commitRepo, repositoryRepo, badgerService)

	// Create handlers for commit and repository operations
	commitHandler := handlers.NewCommitHandler(commitService, repositoryService)
	repositoryHandler := handlers.NewRepositoryHandler(repositoryService, monitorService)

	// Initialize Gin router and configure API routes
	router := gin.Default()
	routes.SetupAPIRoutes(router, commitHandler, repositoryHandler)

	// Define the server port
	PORT := fmt.Sprintf(":%s", cfg.PORT)

	// Start the HTTP server and listen for incoming requests
	if err := router.Run(PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
