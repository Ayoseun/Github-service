package service

import (
	"context"
	"fmt"
	"github-service/internal/core/domain"
	"github-service/internal/ports"
	"github-service/pkg/logger"
	"github-service/pkg/utils"
	"log"
	"time"
)

type MonitorService struct {
	commitService       CommitServiceImpl
	repositoryService   RepositoryServiceImpl
	githubService       ports.GithubImpl
	maxRetryAttempts    int
	initialRetryBackoff time.Duration
}

func NewMonitorService(commitService CommitServiceImpl, repositoryService RepositoryServiceImpl, maxRetryAttempts int, initialRetryBackoff time.Duration, githubService ports.GithubImpl) *MonitorService {
	return &MonitorService{
		commitService:       commitService,
		repositoryService:   repositoryService,
		maxRetryAttempts:    maxRetryAttempts,
		initialRetryBackoff: initialRetryBackoff,
		githubService:       githubService,
	}
}

// MonitorRepository oversees monitoring both repository and commit information for changes.
func (m *MonitorService) MonitorRepository(ctx context.Context, rData domain.RepoData) error {
	retryCount := 0
	for {
		err := m.syncRepositoryAndCommits(ctx, rData)
		if err == nil {
			break
		}

		logger.LogError(err)
		retryCount++
		if retryCount >= m.maxRetryAttempts {
			return err
		}

		backoffDuration := utils.ExponentialBackoff(retryCount, m.initialRetryBackoff)
		time.Sleep(backoffDuration)
	}
	return nil
}

// SyncRepositoryInfo fetches and updates repository information.
func (ms *MonitorService) SyncRepositoryInfo(ctx context.Context, r domain.RepoData) (bool, error) {

	updatedRepository, err := ms.githubService.FetchRepository(ctx, r.Owner, r.RepoName)
	if err != nil {
		return false, err
	}

	ok, err := ms.repositoryService.UpdateInsert(ctx, updatedRepository)
	if !ok {
		return false, err
	}

	return ok, nil
}

// syncRepositoryAndCommits fetches and updates both repository information and commits.
func (m *MonitorService) syncRepositoryAndCommits(ctx context.Context, rData domain.RepoData) error {
	ok, err := m.SyncRepositoryInfo(ctx, rData)

	if err != nil {
		return err
	}
	if ok {
		m.MonitorRepositoryCommits(ctx, rData.RepoName)
	}

	return nil
}

func (m *MonitorService) MonitorRepositoryCommits(ctx context.Context, repositoryName string, startAt ...time.Time) error {
	// Retrieve the last saved commit for the repository
	lastCommit, err := m.commitService.LastCommit(ctx, repositoryName)

	if err != nil { // Handle DB error, except for no rows (no last commit case)
		return fmt.Errorf("could not get last saved commit: %w", err)
	}

	// Get the repository owner and name
	r, err := m.repositoryService.GetRepository(ctx, repositoryName)
	if err != nil {
		return fmt.Errorf("could not get repository owner and name: %w", err)
	}
	fmt.Println(r.CreatedAt)
	// Handle case when there is no last commit

	// No last commit found, use a default date or repository creation date
	// Fetch commits since the repository creation date
	var since time.Time
	since = r.CreatedAt
	if lastCommit != nil {
		// Convert last commit date to time.Time
		since = lastCommit.CommitDate

	}
	// Asynchronously save commits from the last commit date (or repository creation date) to now
	go func() {
		_, err := m.commitService.SaveCommits(ctx, r.Owner, r.Name, since)
		if err != nil {
			log.Printf("error saving commits: %v", err)
		}
	}()

	return nil
}

func (m *MonitorService) AddRepositoryCommitsToMonitor(ctx context.Context, rData domain.RepoData, startAt time.Time) error {
	// Retrieve the last saved commit for the repository
	lastCommit, err := m.commitService.LastCommit(ctx, rData.RepoName)
	if err != nil {
		return fmt.Errorf("could not get last saved commit: %w", err)
	}

	// Initialize the 'since' time variable
	var since time.Time
	if lastCommit != nil {
		// Use the last commit date if it exists
		since = lastCommit.CommitDate
	} else {
		// Fetch repository info to get creation date if no last commit exists
		repo, err := m.repositoryService.GetRepository(ctx, rData.RepoName)
		if err != nil {
			// Attempt to sync repository info if fetching fails
			if _, err := m.SyncRepositoryInfo(ctx, rData); err != nil {
				return err
			}
			// Try fetching repository info again
			repo, err = m.repositoryService.GetRepository(ctx, rData.RepoName)
			if err != nil {
				return fmt.Errorf("could not get repository info: %w", err)
			}
		}

		// Use repository creation date if no last commit date is available
		since = repo.CreatedAt
	}

	// Override the 'since' time with the provided 'startAt' time if it's not zero
	if !startAt.IsZero() {
		since = startAt
	}

	// Asynchronously save commits from the 'since' date to now
	go func() {
		if _, err := m.commitService.SaveCommits(ctx, rData.Owner, rData.RepoName, since); err != nil {
			log.Printf("error saving commits: %v", err)
		}
	}()

	return nil
}
