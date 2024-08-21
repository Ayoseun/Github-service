package service

import (
	"github-service/internal/models"     // Importing the models package from the project's internal directory
	"github-service/internal/repository" // Importing the repository package from the project's internal directory
	"github-service/pkg/github"          // Importing the GitHub package from the project's pkg directory
	"gorm.io/gorm"                       // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// RepositoryService fetches the repository data from GitHub and saves it to the database
func RepositoryService(repo string, db *gorm.DB) (*models.Repository, error) {
	// Create a new instance of the Repository
	rRepo := repository.NewRepository(db)

	// Fetch the repository data from GitHub
	r, err := github.FetchRepositoryData(repo)
	if err != nil {
		// Handle the error, e.g., log it or return the error
		return nil, err
	}

	// Save the repository data to the database
	err = rRepo.SaveRepository(r)
	if err != nil {
		// Handle the error, e.g., log it or return the error
		return nil, err
	}

	// Return the fetched repository data
	return r, nil
}
