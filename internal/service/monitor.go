package service

import (
	"gorm.io/gorm"
	"time"
)

func StartDataFetchingTask(db *gorm.DB, repositoryName string) {
	fetchTicker := time.NewTicker(1 * time.Minute) // Fetch data every minute
	defer fetchTicker.Stop()

	var lastFetchedCommitDate time.Time

	for range fetchTicker.C {
		fetchAndStoreData(db, repositoryName, &lastFetchedCommitDate)
	}
}

func fetchAndStoreData(db *gorm.DB, repositoryName string, lastFetchedCommitDate *time.Time) {
	// Fetch repository commits
	newCommits := CommitsService(repositoryName, db, *lastFetchedCommitDate)

	if len(newCommits) > 0 {
		*lastFetchedCommitDate = newCommits[len(newCommits)-1].Commit.Author.Date
	}

	// Fetch repository data
	RepositoryService(repositoryName, db)
}
