package database

import (
	"github-service/internal/config"
	"github-service/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database using the provided configuration
func Connect(config config.DatabaseConfig) (*gorm.DB, error) {
	// Open a new database connection using the provided data source name (DSN)
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the database schema by automatically creating the necessary tables
	// based on the provided model structures
	db.AutoMigrate(&models.SavedCommit{}, &models.Repository{})

	return db, nil
}
