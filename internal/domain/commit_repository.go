package domain

import "github-service/internal/domain/models"

// CommitRepository defines the interface for commit data operations
type CommitRepository interface {
	SaveCommit(commit *models.SavedCommit) error
	GetCommits(repositoryURL string, page, limit int) ([]models.SavedCommit, error)
	GetTotalCommits(repositoryURL string) (int64, error)
}
