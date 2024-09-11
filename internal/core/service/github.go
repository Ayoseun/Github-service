package service

import (
	"context"
	"fmt"
	"time"

	"github-service/config"
	"github-service/internal/adapters/github"
	"github-service/internal/core/domain"
	"github-service/internal/ports"

	"github-service/pkg/logger"
)

// githubService provides operations for managing commits and config injection
type githubService struct {
	cfg    *config.Config
	client *github.GithubClient
	ctx    context.Context
}

// NewGihubService creates a new instance of CommitService
func NewGithubService(cfg *config.Config, ctx context.Context, client *github.GithubClient) ports.GithubImpl {
	return &githubService{cfg: cfg, ctx: ctx, client: client}
}

// FetchAndSaveCommits fetches commits from GitHub and saves them to the database
func (s *githubService) FetchCommit(ctx context.Context, owner, repo string, since time.Time) ([]domain.Commit, error) {
	commits, err := s.client.FetchRepositoryCommits(ctx, owner, repo, since)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	domainCommits := convertToDomainCommits(commits, repo)

	return domainCommits, nil
}

// FetchAndSaveCommits fetches commits from GitHub and saves them to the database
func (s *githubService) FetchRepository(ctx context.Context, owner, repoName string) (*domain.Repository, error) {
	apiRepo, err := s.client.FetchRepositoryMetaData(ctx, owner, repoName)
	if err != nil {
		logger.LogError(err)
		return &domain.Repository{}, err
	}
	repositoryMetadata := &domain.Repository{
		Owner:            owner,
		Name:             apiRepo.Name,
		Description:      apiRepo.Description,
		URL:              apiRepo.URL,
		Language:         apiRepo.Language,
		ForksCount:       apiRepo.ForksCount,
		StarsGazersCount: apiRepo.StarsGazersCount,
		OpenIssuesCount:  apiRepo.OpenIssuesCount,
		WatchersCount:    apiRepo.WatchersCount,
		CreatedAt:        apiRepo.CreatedAt,
		UpdatedAt:        apiRepo.UpdatedAt,
	}
	logger.LogInfo(fmt.Sprintf("Repository fetched: %s/%s", owner, repoName))
	return repositoryMetadata, nil
}

// convertToDomainCommits converts API commits to domain commits.
func convertToDomainCommits(apiCommits []github.Commit, repo string) []domain.Commit {
	domainCommits := make([]domain.Commit, len(apiCommits))
	for i, commit := range apiCommits {
		domainCommits[i] = domain.Commit{
			Hash:       commit.SHA,
			Message:    commit.Commit.Message,
			Author:     commit.Commit.Committer.Name,
			Email:      commit.Commit.Committer.Email,
			CommitDate: commit.Commit.Committer.Date,
			URL:        commit.Commit.URL,
			Repository: repo,
		}
	}
	return domainCommits
}
