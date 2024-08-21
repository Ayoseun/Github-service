package repository_test

import (
	"github-service/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Automatically migrate the SavedCommit model to create the table
	db.AutoMigrate(&models.SavedCommit{})

	return db
}
