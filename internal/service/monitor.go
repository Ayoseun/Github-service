package service

import (
	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
	"time"         // Importing the time package from the standard library
)

// StartDataFetchingTask starts a background task to periodically fetch and store repository data
func StartDataFetchingTask(db *gorm.DB, repositoryName string) {
	fetchTicker := time.NewTicker(1 * time.Hour) // Set the fetch interval to 1 minute
	defer fetchTicker.Stop()                     // Stop the ticker when the function returns

	var lastFetchedCommitDate time.Time // Initialize the last fetched commit date

	for range fetchTicker.C { // Run the fetch and store loop on the ticker interval
		fetchAndStoreData(db, repositoryName, &lastFetchedCommitDate)
	}
}

// fetchAndStoreData fetches the repository data and stores it in the database
func fetchAndStoreData(db *gorm.DB, repositoryName string, lastFetchedCommitDate *time.Time) {
	// Fetch and store repository data
	RepositoryService(repositoryName, db)
	// Fetch repository commits
	newCommits, err := CommitsService(repositoryName, db, *lastFetchedCommitDate)
	if err != nil {
		panic("Let it panic in this test scenario, because something is definitely wrong")
	}
	if len(newCommits) > 0 {
		// Update the last fetched commit date to the date of the most recent commit
		*lastFetchedCommitDate = newCommits[len(newCommits)-1].Commit.Author.Date
	}

}
