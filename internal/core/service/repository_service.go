package service

import (
	"context"
	"fmt"
	"github-service/config"
	"github-service/pkg/logger"

	"github-service/internal/core/domain"
	"github-service/internal/ports"
)

type RepositoryServiceImpl interface {
	FetchAndSaveRepository(ctx context.Context, rData domain.RepoData) (*domain.Repository, error)
	GetRepository(ctx context.Context, repositoryName string) (domain.Repository, error)
	UpdateInsert(ctx context.Context, d *domain.Repository) (bool, error)
	GetTopNCommitAuthors(ctx context.Context, repositoryName string, n, page, limit int) (domain.TopAuthorsCount, error)
	DeleteARepository(ctx context.Context, owner, repositoryName string) (bool, error)
}

// RepositoryService provides operations for managing repository and cfg injection
type RepositoryService struct {
	postgresRepo  ports.PostgresRepository
	commitService CommitService
	badgerService ports.BadgerImpl
	cfg           *config.Config
	githubService ports.GithubImpl
}

// NewRepositoryService creates a new instance of RepositoryService
func NewRepositoryService(postgresRepo ports.PostgresRepository, commitService CommitService, cfg *config.Config, bs ports.BadgerImpl, githubService ports.GithubImpl) *RepositoryService {
	return &RepositoryService{postgresRepo: postgresRepo, commitService: commitService, cfg: cfg, badgerService: bs, githubService: githubService}
}

// RepositoryService fetches the repository data from GitHub and saves it to the database
func (rs *RepositoryService) FetchAndSaveRepository(ctx context.Context, rData domain.RepoData) (*domain.Repository, error) {
	if err := rs.badgerService.UpdateRepoArray("repos", rData); err != nil {
		logger.LogError(err)
	}
	// Fetch the repository data from GitHub
	r, err := rs.githubService.FetchRepository(ctx, rData.Owner, rData.RepoName)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	ok, err := rs.UpdateInsert(ctx, r)
	if !ok {
		logger.LogError(err)
		return nil, err
	}

	// Return the fetched repository data
	return r, nil
}

func (rs *RepositoryService) UpdateInsert(ctx context.Context, d *domain.Repository) (bool, error) {

	if err := rs.postgresRepo.SaveRepository(ctx, d); err != nil {
		return false, err
	}
	logger.LogInfo(fmt.Sprintf("Saved repository,%s", d.Name))
	return true, nil
}

// GetRepository gets a saved repository by the repository name
func (rs *RepositoryService) GetRepository(ctx context.Context, repositoryName string) (domain.Repository, error) {
	data, err := rs.postgresRepo.GetRepositoryByName(ctx, repositoryName)

	return data, err
}

func (rs *RepositoryService) GetTopNCommitAuthors(ctx context.Context, repoName string, n, page, limit int) (domain.TopAuthorsCount, error) {
	// Calculate how many authors to fetch on this page
	if page*limit > n {
		// Adjust limit to fetch only the remaining authors needed to reach `n`
		limit = n - (page-1)*limit
	}
	// Fetch the authors using pagination
	authors, err := rs.postgresRepo.GetTopNCommitAuthors(ctx, repoName, page, limit)
	return authors, err
}
func (rs *RepositoryService) DeleteARepository(ctx context.Context, owner, repositoryName string) (bool, error) {

	ok, err := rs.postgresRepo.DeleteRepository(ctx, owner, repositoryName)
	if !ok {
		return false, err
	}
	ok, err = rs.commitService.DeleteCommits(ctx, repositoryName)
	if !ok {
		return false, err
	}
	return ok, nil
}
