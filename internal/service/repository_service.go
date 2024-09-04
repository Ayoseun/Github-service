package service

import (
	"github-service/internal/config"
	"github-service/internal/domain"
	"github-service/internal/domain/models"

	"github-service/pkg/github"
)

// RepositoryService provides operations for managing repository and cfg injection
type RepositoryService struct {
	repositoryRepository domain.RepositoryRepository
	cfg                  config.Config
}

// NewRepositoryService creates a new instance of RepositoryService
func NewRepositoryService(repositoryRepository domain.RepositoryRepository, cfg config.Config) *RepositoryService {
	return &RepositoryService{repositoryRepository: repositoryRepository, cfg: cfg}
}

// RepositoryService fetches the repository data from GitHub and saves it to the database
func (s *RepositoryService) FetchAndSaveRepository(owner, repo string) (*models.Repository, error) {

	// Fetch the repository data from GitHub
	r, err := github.FetchRepositoryMetaData(owner, repo, s.cfg)
	if err != nil {
		// Handle the error, e.g., log it or return the error
		return nil, err
	}
	if err := s.repositoryRepository.SaveRepository(r); err != nil {
		return nil, err
	}

	// Return the fetched repository data
	return r, nil
}

// GetRepository gets a saved repository by the repository name
func (s *RepositoryService) GetRepository(repositoryName string) (models.Repository, error) {
	data, err := s.repositoryRepository.GetRepositoryByName(repositoryName)
	return data, err
}

func (s *RepositoryService) GetTopNCommitAuthors(n int, page int, limit int) (models.TopAuthorsCount, error) {
	// Calculate how many authors to fetch on this page
	if page*limit > n {
		// Adjust limit to fetch only the remaining authors needed to reach `n`
		limit = n - (page-1)*limit
	}
	// Fetch the authors using pagination
	authors, err := s.repositoryRepository.GetTopNCommitAuthors(page, limit)
	return authors, err
}
