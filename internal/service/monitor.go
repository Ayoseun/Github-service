package service

import (
	"context"
	"github-service/internal/config"
	"log"
	"time"
)

type CommitMonitorImpl struct {
	commitService     *CommitService
	repositoryService *RepositoryService
}

func NewCommitMonitor(commitService *CommitService, repositoryService *RepositoryService) *CommitMonitorImpl {
	return &CommitMonitorImpl{
		commitService:     commitService,
		repositoryService: repositoryService,
	}
}

// StartDataFetchingTask starts a background task to periodically fetch and store repository data
func (h *CommitMonitorImpl) StartDataFetchingTask(ctx context.Context, cfg config.Config, repositoryOwner, repositoryName string) {
	// Create a new ticker that triggers every 5 hours
	fetchTicker := time.NewTicker(5 * time.Second)
	defer fetchTicker.Stop() // Ensure the ticker is stopped when the function exits

	var lastFetchedCommitDate time.Time // Variable to store the date of the last fetched commit

	log.Println("Data fetching task started.")
	// Start an infinite loop that runs every time the ticker ticks
	for {
		select {
		case <-ctx.Done():
			// Context has been canceled or timed out
			var now time.Time
			log.Printf("Data fetching task stopped at : %s", now)
			return
		case <-fetchTicker.C:
			// Call the fetchAndStoreData function to fetch new data and store it in the database
			log.Println("Fetching and storing data...")
			if _, err := h.fetchAndStoreData(repositoryOwner, repositoryName, &lastFetchedCommitDate); err != nil {
				log.Printf("Error fetching and storing data: %v", err)
			} else {
				log.Println("Data successfully fetched and stored.")
			}
		}
	}
}

// SeedDB seeds the database with initial data by fetching repository commits
func (h *CommitMonitorImpl) SeedDB(repositoryOwner, repositoryName string, beginFetchCommitDate time.Time) {
	// Call the fetchAndStoreData function to fetch and store data, starting from the specified date
	h.fetchAndStoreData(repositoryOwner, repositoryName, &beginFetchCommitDate)
}

// fetchAndStoreData fetches the repository data and stores it in the database
func (h *CommitMonitorImpl) fetchAndStoreData(repositoryOwner, repositoryName string, lastFetchedCommitDate *time.Time) (any, error) {
	// Fetch repository data using the RepositoryService and store it in the database
	_, err := h.repositoryService.FetchAndSaveRepository(repositoryOwner, repositoryName)
	if err != nil {
		log.Printf("Failed to fetch repository data for %s/%s: %v", repositoryOwner, repositoryName, err)
		return nil, err
	}

	// Fetch new commits from the repository, starting from the last fetched commit date
	newCommits, err := h.commitService.FetchAndSaveCommits(repositoryOwner, repositoryName, *lastFetchedCommitDate)
	if err != nil {
		log.Printf("Failed to fetch commits for %s/%s: %v", repositoryOwner, repositoryName, err)
		return nil, err
	}

	// If there are new commits, update the last fetched commit date
	if len(newCommits) > 0 {
		*lastFetchedCommitDate = newCommits[len(newCommits)-1].Commit.Author.Date
	}

	return nil, err
}
