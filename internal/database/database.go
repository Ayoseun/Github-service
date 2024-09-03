package database

import (
	"fmt"
	"github-service/internal/config"
	"github-service/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// Connect establishes a connection to the database using the provided configuration
// It will retry up to 3 times before returning an error
func Connect(config config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	maxRetries := 3
	retryDelay := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			fmt.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...\n",
				i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Migrate the database schema by automatically creating the necessary tables
	// based on the provided model structures
	err = db.AutoMigrate(&models.SavedCommit{}, &models.Repository{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %v", err)
	}

	return db, nil
}
