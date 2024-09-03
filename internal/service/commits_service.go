package service

import (
	"github-service/internal/config"
	"github-service/internal/domain"
	"github-service/internal/domain/models"
	"github-service/pkg/github"
	"time"
)

// CommitService provides operations for managing commits and config injection
type CommitService struct {
	commitRepository domain.CommitRepository
	cfg              config.Config
}

// NewCommitService creates a new instance of CommitService
func NewCommitService(commitRepository domain.CommitRepository, cfg config.Config) *CommitService {
	return &CommitService{commitRepository: commitRepository, cfg: cfg}
}

// FetchAndSaveCommits fetches commits from GitHub and saves them to the database
func (s *CommitService) FetchAndSaveCommits(owner, repo string, lastFetchedCommitDate time.Time) ([]models.Commit, error) {
	commits, err := github.FetchRepositoryCommits(owner, repo, s.cfg)
	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		if commit.Commit.Author.Date.After(lastFetchedCommitDate) {
			commitToSave := &models.SavedCommit{
				URL:        commit.URL,
				Message:    commit.Commit.Message,
				Author:     commit.Commit.Author.Name,
				Date:       commit.Commit.Author.Date,
				Repository: repo,
			}

			if err := s.commitRepository.SaveCommit(commitToSave); err != nil {
				return nil, err
			}
		}
	}

	return commits, nil
}

// GetPaginatedCommits returns paginated commits from the database
func (s *CommitService) GetPaginatedCommits(repository string, page, limit int) ([]models.SavedCommit, error) {
	return s.commitRepository.GetCommits(repository, page, limit)
}

// GetCommitCount returns the total number of commits for a repository URL
func (s *CommitService) GetCommitCount(repository string) (int64, error) {
	return s.commitRepository.GetTotalCommits(repository)
}
