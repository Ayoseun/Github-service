package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"github-service/internal/core/domain"
	"github-service/internal/ports"

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// RepositoryImpl is a struct that implements the PostgresRepository interface using GORM
type RepositoryImpl struct {
	DB *gorm.DB // Holds a reference to the database connection
}

// NewRepository is a constructor function that creates a new instance of RepositoryImpl
// It returns an error if the provided database connection is nil
func NewRepository(db *gorm.DB) (ports.PostgresRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &RepositoryImpl{DB: db}, nil
}

// SaveRepository saves a repository to the database, creating a new one if it doesn't exist, or updating an existing one
func (r *RepositoryImpl) SaveRepository(ctx context.Context, repository *domain.Repository) error {
	var existingRepo domain.Repository
	err := r.DB.WithContext(ctx).Where("name = ?", repository.Name).First(&existingRepo).Error
	if err != nil {
		// If the repository does not exist, create a new one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.DB.WithContext(ctx).Create(repository).Error
		}
		return err
	}
	// Update the existing repository
	return r.DB.WithContext(ctx).Model(&existingRepo).Updates(repository).Error
}

// GetTopNCommitAuthors retrieves the top N commit authors, with pagination support
func (r *RepositoryImpl) GetTopNCommitAuthors(ctx context.Context, repositoryName string, page, limit int) (domain.TopAuthorsCount, error) {
	var authors domain.TopAuthorsCount

	err := r.DB.WithContext(ctx).
		Model(&domain.Commit{}).
		Where("repository = ?", repositoryName).
		Select("author, count(author) as count").
		Group("author").
		Order("count DESC, author ASC"). // Order by count descending, then by author name ascending
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&authors).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.TopAuthorsCount{}, fmt.Errorf("no authors found for repository %s: %w", repositoryName, err)
		}
		return domain.TopAuthorsCount{}, fmt.Errorf("failed to retrieve authors: %w", err)
	}
	return authors, nil
}

// GetRepositoryByName retrieves a repository based on its name
func (r *RepositoryImpl) GetRepositoryByName(ctx context.Context, repositoryName string) (domain.Repository, error) {
	var repository domain.Repository
	err := r.DB.WithContext(ctx).Where("name = ?", repositoryName).First(&repository).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Repository{}, fmt.Errorf("repository %s not found", repositoryName)
	}
	return repository, err
}

// DeleteRepository deletes a repository with the given name
func (r *RepositoryImpl) DeleteRepository(ctx context.Context, owner, repositoryName string) (bool, error) {
	result := r.DB.WithContext(ctx).Where("name = ? AND owner = ?", repositoryName, owner).Delete(&domain.Repository{})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil // No repository found with the given name
	}
	return true, nil // Repository successfully deleted
}
