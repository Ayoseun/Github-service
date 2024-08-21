package service

import (
	"github-service/internal/models"     // Importing the models package from the project's internal directory
	"github-service/internal/repository" // Importing the repository package from the project's internal directory
	"github-service/pkg/github"          // Importing the GitHub package from the project's pkg directory
	"time"                               // Importing the time package from the standard library

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// CommitsService fetches the commits for a given repository and saves the new commits to the database
func CommitsService(repo string, db *gorm.DB, lastFetchedCommitDate time.Time) ([]models.Commit, error) {
	// Create a new instance of the CommitRepository
	commitRepo := repository.NewCommitRepository(db)

	// Fetch the repository commits from GitHub
	commits, err := github.FetchRepositoryCommits(repo)
	if err != nil {
		return []models.Commit{}, err
	}

	// Iterate through the fetched commits and save the new ones to the database
	for _, commit := range commits {
		if commit.Commit.Author.Date.After(lastFetchedCommitDate) {
			// Create a new SavedCommit model with the commit data
			savedCommit := &models.SavedCommit{
				Message: commit.Commit.Message,
				Author:  commit.Commit.Author.Name,
				Date:    commit.Commit.Author.Date,
				URL:     commit.HTMLURL,
			}

			// Save the new commit to the database
			commitRepo.SaveCommits(savedCommit)
		}
	}

	// Return the fetched commits
	return commits, err
}
