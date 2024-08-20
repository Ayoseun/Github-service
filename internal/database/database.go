package database

import (
	"github-service/internal/config"
	"github-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config config.DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.SavedCommit{}, &models.Repository{})

	return db
}
