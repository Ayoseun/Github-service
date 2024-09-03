package domain

import "github-service/internal/domain/models"

// Repository defines the interface for commit data operations
type RepositoryRepository interface {
	SaveRepository(repo *models.Repository) error
	GetTopNCommitAuthors(page, limit int) (models.TopAuthorsCount, error)
	GetRepositoryByName(repository string) (models.Repository, error)
}
