package service

import (
	"context"
	"fmt"

	"github-service/config"
	"github-service/internal/core/domain"
	"github-service/internal/ports"
	"time"

	"github-service/pkg/logger"
)

type CommitServiceImpl interface {
	SaveCommits(ctx context.Context, owner, repoName string, since time.Time, until ...time.Time) ([]domain.Commit, error)
	GetPaginatedCommits(ctx context.Context, repositoryName string, page, limit int) ([]domain.Commit, error)
	GetCommitCount(ctx context.Context, repositoryName string) (int64, error)
	DeleteCommits(ctx context.Context, repositoryName string) (bool, error)
	LastCommit(ctx context.Context, repositoryName string) (*domain.Commit, error)
}

// CommitService provides operations for managing commits and config injection
type CommitService struct {
	pc            ports.PostgresCommit
	cfg           *config.Config
	githubService ports.GithubImpl
}

// NewCommitService creates a new instance of CommitService
func NewCommitService(postgresCommitRepository ports.PostgresCommit, cfg *config.Config, githubService ports.GithubImpl) *CommitService {
	return &CommitService{pc: postgresCommitRepository, cfg: cfg, githubService: githubService}
}

func (cs *CommitService) SaveCommits(ctx context.Context, owner, repoName string, since time.Time, until ...time.Time) ([]domain.Commit, error) {

	// Fetch commits from GitHub
	commits, err := cs.githubService.FetchCommit(ctx, owner, repoName, since)
	fmt.Println(commits)
	logger.LogInfo("fetching commiting from github")
	if err != nil {
		return nil, err
	}
	for _, commit := range commits {
		if err := cs.pc.SaveCommit(ctx, &commit); err != nil {
			return nil, err
		}
	}
	// var savedCommits []domain.Commit

	// // Determine if "until" was provided
	// if len(until) > 0 {
	// 	// Case: "until" is provided, use both "since" and "until"
	// 	untilDate := until[0] // Use the first value from the "until" slice

	// 	// Loop through the commits and filter based on "since" and "until"
	// 	for _, commit := range commits {

	// 		// Save commits that are after "since" and before "until"
	// 		if commit.CommitDate.After(since) && commit.CommitDate.Before(untilDate) {
	// 			if err := cs.pc.SaveCommit(ctx, &commit); err != nil {
	// 				return nil, err
	// 			}
	// 			savedCommits = append(savedCommits, commit) // Track saved commits
	// 		}
	// 	}
	// } else {
	// 	// Case: "until" is not provided, only use "since"
	// 	// Loop through the commits and filter based on "since" only
	// 	for _, commit := range commits {

	// 		// Save commits that are after "since" (no upper bound)
	// 		if commit.CommitDate.After(since) {
	// 			if err := cs.pc.SaveCommit(ctx, &commit); err != nil {
	// 				return nil, err
	// 			}
	// 			savedCommits = append(savedCommits, commit) // Track saved commits
	// 		}
	// 	}
	// }

	return commits, nil
}

// GetPaginatedCommits returns paginated commits from the database
func (cs *CommitService) GetPaginatedCommits(ctx context.Context, repositoryName string, page, limit int) ([]domain.Commit, error) {
	return cs.pc.GetCommits(ctx, repositoryName, page, limit)
}

// GetCommitCount returns the total number of commits for a repository URL
func (cs *CommitService) GetCommitCount(ctx context.Context, repositoryName string) (int64, error) {
	return cs.pc.GetTotalCommits(ctx, repositoryName)
}

func (cs *CommitService) DeleteCommits(ctx context.Context, repositoryName string) (bool, error) {
	if repositoryName == "" {
		return false, fmt.Errorf("Repository cannot be empty")
	}
	ok, err := cs.pc.DeleteAllCommits(ctx, repositoryName)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (cs *CommitService) LastCommit(ctx context.Context, repositoryName string) (*domain.Commit, error) {
	if repositoryName == "" {
		return &domain.Commit{}, fmt.Errorf("Repository cannot be empty")
	}
	commit, err := cs.pc.GetLastCommitByRepositoryName(ctx, repositoryName)
	if err != nil {
		return &domain.Commit{}, err
	}
	return commit, nil
}
