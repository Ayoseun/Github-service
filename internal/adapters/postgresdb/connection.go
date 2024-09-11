package postgresdb

import (
	"fmt"
	"github-service/config"
	"github-service/internal/core/domain"
	"github-service/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database using the provided configuration
// It will retry up to 3 times before returning an error.
func Connect(cfg config.Config) (*gorm.DB, error) {
	// Build the DSN (Data Source Name) for PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_HOST, cfg.POSTGRES_DB)

	var db *gorm.DB
	var err error
	maxRetries := 3
	retryDelay := 5 * time.Second

	// Retry connecting to the database
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break // Successfully connected
		}

		// Log and retry
		if i < maxRetries-1 {
			logger.LogWarning(fmt.Sprintf("Failed to connect to the database (attempt %d/%d): %v. Retrying in %v...",
				i+1, maxRetries, err, retryDelay))
			time.Sleep(retryDelay)
		}
	}

	// Return error after exhausting retries
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after %d attempts: %w", maxRetries, err)
	}

	logger.LogInfo("Successfully connected to the database.")

	// Automatically migrate the schema (create/update tables based on the provided models)
	err = db.AutoMigrate(&domain.Commit{}, &domain.Repository{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %v", err)
	}

	return db, nil
}
