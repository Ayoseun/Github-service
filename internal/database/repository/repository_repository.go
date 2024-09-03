package repository

import (
	"errors"
	"github-service/internal/domain"
	"github-service/internal/domain/models" // Importing the models package from the project's internal directory

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// RepositoryImpl is a struct that represents a repository for managing repositories and their commit authors
// It implements the RepositoryRepository interface
type RepositoryImpl struct {
	DB *gorm.DB // Holds a reference to the database connection
}

// NewRepository is a constructor function that creates a new instance of RepositoryImpl
// It returns an error if the provided database connection is nil
func NewRepository(db *gorm.DB) (domain.RepositoryRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &RepositoryImpl{DB: db}, nil
}

// SaveRepository saves a repository to the database, creating a new one if it doesn't exist, or updating an existing one
// It checks if the repository already exists based on its name and either creates or updates it accordingly
func (r *RepositoryImpl) SaveRepository(repository *models.Repository) error {
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
	// Update existing repository's data using GORM's auto-update functionality
	return r.DB.Updates(&existingRepo).Error
}

// GetTopNCommitAuthors retrieves the top N commit authors, with pagination support
// It groups authors by name, counts their commits, and orders the results by the count in descending order
func (r *RepositoryImpl) GetTopNCommitAuthors(page, limit int) (models.TopAuthorsCount, error) {
	var authors models.TopAuthorsCount

	err := r.DB.Model(&models.SavedCommit{}).
		Select("author, count(author) as count").
		Group("author").
		Order("count DESC, author ASC"). // Order by count descending, then by author name ascending
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&authors).Error

	return authors, err
}

// GetRepositoryByURL retrieves a repository based on its URL
// It returns the repository model and an error if the query fails
func (r *RepositoryImpl) GetRepositoryByURL(repositoryURL string) (models.Repository, error) {

	var repository models.Repository
	err := r.DB.Where("name = ?", repositoryURL).First(&repository).Error

	return repository, err
}
