package ports

import "github-service/internal/core/domain"

// BadgerImpl defines the interface for interacting with the Badger key-value store
type BadgerImpl interface {
	// SaveRepoArray saves an array of RepoData under a specified repoKey in the Badger store
	// Returns an error if the operation fails
	SaveRepoArray(repoKey string, repoDataArray []domain.RepoData) error

	// GetRepoArray retrieves the array of RepoData stored under the specified repoKey
	// Returns an empty slice and error if the key does not exist or retrieval fails
	GetRepoArray(repoKey string) ([]domain.RepoData, error)

	// UpdateRepoArray updates a specific RepoData within the array stored under repoKey
	// It fetches the array, modifies the target RepoData, and saves it back to the Badger store
	// Returns an error if the update operation fails
	UpdateRepoArray(repoKey string, updatedRepo domain.RepoData) error
}
