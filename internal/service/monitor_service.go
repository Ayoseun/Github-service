package service

import (
	"context"
	"github-service/internal/domain/models"
	"log"
	"sync"
	"time"
)

type CommitMonitor struct {
	commitService     CommitServiceInterface
	repositoryService RepositoryServiceInterface
	mu                sync.Mutex // Protects the repositories map
	repositories      map[string]context.CancelFunc
	ctx               context.Context
}

type CommitServiceInterface interface {
	FetchAndSaveCommits(owner, repo string, since time.Time) ([]models.Commit, error)
	FetchCommitsInRange(owner, repo string, from, to time.Time) ([]models.Commit, error)
}

type RepositoryServiceInterface interface {
	FetchAndSaveRepository(owner, repo string) (*models.Repository, error)
}

func NewCommitMonitor(ctx context.Context, commitService CommitServiceInterface, repositoryService RepositoryServiceInterface) *CommitMonitor {
	return &CommitMonitor{
		commitService:     commitService,
		repositoryService: repositoryService,
		repositories:      make(map[string]context.CancelFunc),
		ctx:               ctx,
	}
}

// StartDataFetchingTask starts a background task to periodically fetch and store repository data
func (cm *CommitMonitor) StartDataFetchingTask(ctx context.Context, owner, repo string) {
	const fetchInterval = 5 * time.Second

	fetchTicker := time.NewTicker(fetchInterval)
	defer fetchTicker.Stop()

	var lastFetchedCommitDate time.Time

	log.Printf("Data fetching task started for %s/%s.", owner, repo)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Data fetching task stopped for %s/%s.", owner, repo)
			return
		case <-fetchTicker.C:
			log.Printf("Fetching and storing data for %s/%s...", owner, repo)
			if err := cm.fetchAndStoreData(owner, repo, &lastFetchedCommitDate); err != nil {
				log.Printf("Error fetching and storing data for %s/%s: %v", owner, repo, err)
			} else {
				log.Printf("Data successfully fetched and stored for %s/%s.", owner, repo)
			}
		}
	}
}

// AddRepository adds a new repository and starts monitoring it. Optionally, a range of commit dates can be provided.
func (cm *CommitMonitor) AddRepository(owner, repo string, since time.Time, dateRange ...time.Time) (bool, error) {
	if len(dateRange) == 2 {
		cm.fetchAndStoreData(owner, repo, &since, dateRange[0], dateRange[1])
	} else {
		cm.fetchAndStoreData(owner, repo, &since)
	}

	ctx, cancel := context.WithCancel(cm.ctx)

	cm.mu.Lock()
	cm.repositories[owner+"/"+repo] = cancel
	cm.mu.Unlock()

	go cm.StartDataFetchingTask(ctx, owner, repo)
	return true, nil
}

// RemoveRepository stops monitoring a repository by canceling its background task
func (cm *CommitMonitor) RemoveRepository(owner, repo string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	repoKey := owner + "/" + repo

	if cancel, exists := cm.repositories[repoKey]; exists {
		cancel()
		delete(cm.repositories, repoKey)
		log.Printf("Repository %s/%s removed and monitoring stopped.", owner, repo)
	} else {
		log.Printf("Repository %s/%s not found.", owner, repo)
	}
}

// fetchAndStoreData fetches the repository data and stores it in the database
func (cm *CommitMonitor) fetchAndStoreData(owner, repo string, since *time.Time, dateRange ...time.Time) error {
	_, err := cm.repositoryService.FetchAndSaveRepository(owner, repo)
	if err != nil {
		return err
	}

	var newCommits []models.Commit
	if len(dateRange) == 2 {
		newCommits, err = cm.commitService.FetchCommitsInRange(owner, repo, dateRange[0], dateRange[1])
	} else {
		newCommits, err = cm.commitService.FetchAndSaveCommits(owner, repo, *since)
	}
	if err != nil {
		return err
	}

	if len(newCommits) > 0 {
		*since = newCommits[len(newCommits)-1].Commit.Author.Date
	}

	return nil
}
