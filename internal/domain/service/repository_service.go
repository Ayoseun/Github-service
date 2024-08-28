package service

import (
	"github-service/internal/database/repository"
	"github-service/internal/domain/models" // Importing the models package from the project's internal directory
	// Importing the repository package from the project's internal directory
	"github-service/pkg/github" // Importing the GitHub package from the project's pkg directory

	"gorm.io/gorm" // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// RepositoryService fetches the repository data from GitHub and saves it to the database
func RepositoryService(repo string, db *gorm.DB) (*models.Repository, error) {

	// Fetch the repository data from GitHub
	r, err := github.FetchRepositoryData(repo)
	if err != nil {
		// Handle the error, e.g., log it or return the error
		return nil, err
	}
	// Create a new instance of the Repository
	dbInstance := repository.NewRepository(db)

	// Save the repository data to the database
	err = dbInstance.SaveRepository(r)
	if err != nil {
		// Handle the error, e.g., log it or return the error
		return nil, err
	}

	// Return the fetched repository data
	return r, nil
}
