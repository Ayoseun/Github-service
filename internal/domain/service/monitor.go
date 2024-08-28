package service

import (
	"fmt"
	"time" // Importing the time package from the standard library

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// StartDataFetchingTask starts a background task to periodically fetch and store repository data
// It takes in a GORM database instance and the name of the repository to fetch data from
func StartDataFetchingTask(db *gorm.DB, repositoryName string) {
	// Create a new ticker that triggers every hour
	fetchTicker := time.NewTicker(5 * time.Second)
	defer fetchTicker.Stop() // Ensure the ticker is stopped when the function exits

	var lastFetchedCommitDate time.Time // Variable to store the date of the last fetched commit

	// Start an infinite loop that runs every time the ticker ticks (every hour)
	for range fetchTicker.C {
		// Call the fetchAndStoreData function to fetch new data and store it in the database
		fetchAndStoreData(db, repositoryName, &lastFetchedCommitDate)
	}
}

// SeedDB seeds the database with initial data by fetching repository commits
// It takes in a GORM database instance, the repository name, and the date from which to begin fetching commits
func SeedDB(db *gorm.DB, repositoryName string, beginFetchCommitDate time.Time) {
	// Call the fetchAndStoreData function to fetch and store data, starting from the specified date
	fetchAndStoreData(db, repositoryName, &beginFetchCommitDate)
}

// fetchAndStoreData fetches the repository data and stores it in the database
// It takes in a GORM database instance, the repository name, and a pointer to the last fetched commit date
func fetchAndStoreData(db *gorm.DB, repositoryName string, lastFetchedCommitDate *time.Time) {
	// Fetch repository data using the RepositoryService function and store it in the database
	RepositoryService(repositoryName, db)

	// Fetch new commits from the repository, starting from the last fetched commit date
	newCommits, err := CommitsService(repositoryName, db, *lastFetchedCommitDate)
	if err != nil {
		// Log an error message if there's an issue connecting to the database
		fmt.Errorf("failed to connect to database: %v", err)
	}

	// If there are new commits, update the last fetched commit date
	if len(newCommits) > 0 {
		// Set the last fetched commit date to the date of the most recent commit
		*lastFetchedCommitDate = newCommits[len(newCommits)-1].Commit.Author.Date
	}
}
