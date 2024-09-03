package repository

import (
	"errors"
	"github-service/internal/domain"
	"github-service/internal/domain/models"

	"gorm.io/gorm"
)

// CommitRepositoryImpl implements the CommitRepository interface using GORM
type CommitRepositoryImpl struct {
	DB *gorm.DB
}

// NewCommitRepository creates a new instance of CommitRepositoryImpl
func NewCommitRepository(db *gorm.DB) (domain.CommitRepository, error) {

	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &CommitRepositoryImpl{DB: db}, nil
}

// SaveCommit saves a commit to the database
func (r *CommitRepositoryImpl) SaveCommit(commit *models.SavedCommit) error {

	return r.DB.Create(commit).Error
}

// GetCommits retrieves a list of commits based on the repository name, page, and limit
func (r *CommitRepositoryImpl) GetCommits(repository string, page, limit int) ([]models.SavedCommit, error) {
	var commits []models.SavedCommit
	err := r.DB.Where("repository = ?", repository).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&commits).Error
	return commits, err
}

// GetTotalCommits retrieves the total number of commits for the provided repository name
func (r *CommitRepositoryImpl) GetTotalCommits(repository string) (int64, error) {
	var totalCommits int64
	err := r.DB.Where("repository = ?", repository).
		Model(&models.SavedCommit{}).
		Count(&totalCommits).Error
	return totalCommits, err
}
