package repository

import (
	"github-service/internal/models" // Importing the models package from the project's internal directory
	"gorm.io/gorm"                   // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// Repository is a struct that represents a repository for managing repositories and their commit authors
type Repository struct {
	DB *gorm.DB // Holds a reference to the database connection
}

// NewRepository is a constructor function that creates a new instance of the Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db, // Initializing the DB field with the provided GORM database instance
	}
}

// SaveRepository saves a repository to the database, creating a new one if it doesn't exist, or updating an existing one
func (r *Repository) SaveRepository(repository *models.Repository) error {
	// Check if the repository already exists
	var existingRepo models.Repository
	result := r.DB.Where("name = ?", repository.Name).First(&existingRepo)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Repository doesn't exist, create a new one
			return r.DB.Create(repository).Error
		} else {
			// Other error occurred, return it
			return result.Error
		}
	}

	// Repository exists, update it
	existingRepo.Description = repository.Description
	existingRepo.URL = repository.URL
	existingRepo.Language = repository.Language
	existingRepo.ForksCount = repository.ForksCount
	existingRepo.StarsCount = repository.StarsCount
	existingRepo.OpenIssues = repository.OpenIssues
	existingRepo.Watchers = repository.Watchers
	existingRepo.CreatedAt = repository.CreatedAt
	existingRepo.UpdatedAt = repository.UpdatedAt
	existingRepo.SubscribersCount = repository.SubscribersCount
	return r.DB.Save(&existingRepo).Error
}

// GetTopNCommitAuthors retrieves the top N commit authors, with pagination support
func (r *Repository) GetTopNCommitAuthors(n, page, limit int) ([]struct {
	Author string `json:"author"`
	Count  int    `json:"count"`
}, error) {
	var authors []struct {
		Author string `json:"author"`
		Count  int    `json:"count"`
	}

	err := r.DB.Model(&models.SavedCommit{}).
		Select("author, count(author) as count").
		Group("author").
		Order("count desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&authors).Error

	return authors, err
}
