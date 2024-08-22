package repository

import (
	"fmt"                                   // Standard library package for formatting and printing
	"github-service/internal/domain/models" // Importing the models package from the project's internal directory
	"gorm.io/gorm"                          // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// CommitRepository is a struct that represents a repository for managing commits
type CommitRepository struct {
	DB *gorm.DB // Holds a reference to the database connection
}

// NewCommitRepository is a constructor function that creates a new instance of the CommitRepository
func NewCommitRepository(db *gorm.DB) *CommitRepository {
	return &CommitRepository{
		DB: db, // Initializing the DB field with the provided GORM database instance
	}
}

// SaveCommits saves a commit to the database
func (r *CommitRepository) SaveCommits(commit *models.SavedCommit) error {
	return r.DB.Create(commit).Error // Using the GORM Create method to save the commit to the database
}

// GetCommits retrieves a list of commits based on the provided repository URL, page, and limit
func (r *CommitRepository) GetCommits(repositoryURL string, page, limit int) ([]models.SavedCommit, error) {
	var commits []models.SavedCommit // Initializing a slice to hold the retrieved commits

	// Using the GORM Where, Limit, and Offset methods to filter and paginate the commits
	err := r.DB.Where("url LIKE ?", fmt.Sprintf("%%%s%%", repositoryURL)).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&commits).Error

	return commits, err
}

// GetTotalCommits retrieves the total number of commits for the provided repository URL
func (r *CommitRepository) GetTotalCommits(repositoryURL string) (int64, error) {
	var totalCommits int64 // Initializing a variable to hold the total number of commits

	// Using the GORM Where and Count methods to retrieve the total number of commits
	err := r.DB.Where("url LIKE ?", fmt.Sprintf("%%%s%%", repositoryURL)).
		Model(&models.SavedCommit{}).
		Count(&totalCommits).Error

	return totalCommits, err
}
