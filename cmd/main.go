package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	gracefulShutdown(router, PORT)

}

func gracefulShutdown(router *gin.Engine, port string) {
	// Create a channel to listen for OS signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	// Create a server instance with a timeout
	srv := &http.Server{
		Addr:    port,
		Handler: router,
		// Optional: configure timeouts as needed
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("Shutting down server...")

	// Create a timeout context for shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exiting")
}
