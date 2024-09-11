package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"github-service/internal/core/domain"
	"github-service/internal/ports"
	"github-service/pkg/logger"

	"gorm.io/gorm"
)

// CommitRepositoryImpl implements the CommitRepository interface using GORM for database operations.
type CommitRepositoryImpl struct {
	DB *gorm.DB
}

// NewCommitRepository creates a new instance of CommitRepositoryImpl and initializes the database.
// It returns an error if the provided database connection is nil.
func NewCommitRepository(db *gorm.DB) (ports.PostgresCommit, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &CommitRepositoryImpl{DB: db}, nil
}

// SaveCommit saves a commit to the database.
// It returns an error if the save operation fails.
func (c *CommitRepositoryImpl) SaveCommit(ctx context.Context, commit *domain.Commit) error {
	if err := c.DB.WithContext(ctx).Create(commit).Error; err != nil {
		logger.LogWarning(fmt.Sprintf("Failed to save commit for repository %s: %v", commit.Repository, err))
		return err
	}
	logger.LogInfo(fmt.Sprintf("Successfully saved commit for repository %s", commit.Repository))
	return nil
}

// GetCommits retrieves a list of commits based on the repository name, page, and limit.
// The page and limit parameters control pagination.
// It returns a slice of Commit and an error if the query fails.
func (c *CommitRepositoryImpl) GetCommits(ctx context.Context, repositoryName string, page, limit int) ([]domain.Commit, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit must be greater than 0")
	}

	var commits []domain.Commit
	err := c.DB.WithContext(ctx).
		Where("repository = ?", repositoryName).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&commits).Error

	if err != nil {
		logger.LogWarning(fmt.Sprintf("Failed to retrieve commits for repository %s: %v", repositoryName, err))
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("Successfully retrieved commits for repository %s", repositoryName))
	return commits, nil
}

// GetTotalCommits retrieves the total number of commits for the provided repository name.
// It returns the total count and an error if the query fails.
func (c *CommitRepositoryImpl) GetTotalCommits(ctx context.Context, repositoryName string) (int64, error) {
	var totalCommits int64
	err := c.DB.WithContext(ctx).
		Where("repository = ?", repositoryName).
		Model(&domain.Commit{}).
		Count(&totalCommits).Error

	if err != nil {
		logger.LogWarning(fmt.Sprintf("Failed to count commits for repository %s: %v", repositoryName, err))
		return 0, err
	}

	logger.LogInfo(fmt.Sprintf("Successfully retrieved total commits count for repository %s", repositoryName))
	return totalCommits, nil
}

// DeleteAllCommits deletes all commits for the given repository name.
// It returns a boolean indicating success and an error if the delete operation fails.
func (c *CommitRepositoryImpl) DeleteAllCommits(ctx context.Context, repositoryName string) (bool, error) {
	if repositoryName == "" {
		return false, errors.New("repository name must not be empty")
	}

	result := c.DB.WithContext(ctx).Where("repository = ?", repositoryName).Delete(&domain.Commit{})
	if result.Error != nil {
		logger.LogWarning(fmt.Sprintf("Failed to delete commits for repository %s: %v", repositoryName, result.Error))
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		logger.LogWarning(fmt.Sprintf("No commits found for repository %s to delete", repositoryName))
		return false, nil
	}

	logger.LogInfo(fmt.Sprintf("Successfully deleted all commits for repository %s", repositoryName))
	return true, nil
}

// GetLastCommitByRepositoryName retrieves the latest commit for a given repository name.
// It returns nil if no commit is found, or an error if the query fails.
func (c *CommitRepositoryImpl) GetLastCommitByRepositoryName(ctx context.Context, repoName string) (*domain.Commit, error) {
	var commit domain.Commit
	err := c.DB.WithContext(ctx).
		Where("repository = ?", repoName).
		Order("commit_date DESC").
		First(&commit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning(fmt.Sprintf("No commit found for repository %s", repoName))
			return nil, nil
		}
		logger.LogWarning(fmt.Sprintf("Failed to retrieve latest commit for repository %s: %v", repoName, err))
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("Successfully retrieved latest commit for repository %s", repoName))
	return &commit, nil
}
