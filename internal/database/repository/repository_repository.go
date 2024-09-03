package repository

import (
	"errors"
	"github-service/internal/domain"
	"github-service/internal/domain/models" // Importing the models package from the project's internal directory

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// Repository is a struct that represents a repository for managing repositories and their commit authors
type RepositoryImpl struct {
	DB *gorm.DB // Holds a reference to the database connection
}

// NewRepository is a constructor function that creates a new instance of the Repository
func NewRepository(db *gorm.DB) (domain.RepositoryRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &RepositoryImpl{DB: db}, nil
}

// SaveRepository saves a repository to the database, creating a new one if it doesn't exist, or updating an existing one
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
	// Update existing Repository date using gorm auto update
	return r.DB.Updates(&existingRepo).Error
}

// GetTopNCommitAuthors retrieves the top N commit authors, with pagination support
func (r *RepositoryImpl) GetTopNCommitAuthors(page, limit int) (models.TopAuthorsCount, error) {
	var authors models.TopAuthorsCount

	err := r.DB.Model(&models.SavedCommit{}).
		Select("author, count(author) as count").
		Group("author").
		Order("count desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&authors).Error

	return authors, err
}

// GetCommits retrieves a list of commits based on the repository URL, page, and limit
func (r *RepositoryImpl) GetRepositoryByURL(repositoryURL string) (models.Repository, error) {

	var repository models.Repository
	err := r.DB.Where("name = ?", repositoryURL).First(&repository).Error

	return repository, err
}
