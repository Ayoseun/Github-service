package badger

import (
	"encoding/json"
	"fmt"
	"github-service/internal/core/domain"
	"log"

	"github.com/dgraph-io/badger/v3"
)

// BadgerRepository provides methods to interact with a Badger DB instance
type BadgerRepository struct {
	db *badger.DB
}

// NewBadgerRepository initializes a new Badger repository.
// It opens the Badger database at the provided dbPath and returns a repository instance.
func NewBadgerRepository(dbPath string) (*BadgerRepository, error) {
	opts := badger.DefaultOptions(dbPath).WithLogger(nil) // Disable default logging for cleaner output
	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("Error opening Badger database at path %s: %v", dbPath, err)
		return nil, err
	}
	log.Println("Badger database opened successfully")
	return &BadgerRepository{db: db}, nil
}

// SaveRepoArray saves an array of RepoData to Badger using the provided repoKey.
func (b *BadgerRepository) SaveRepoArray(repoKey string, repoDataArray []domain.RepoData) error {
	return b.db.Update(func(txn *badger.Txn) error {
		// Marshal the array of RepoData into JSON
		data, err := json.Marshal(repoDataArray)
		if err != nil {
			log.Printf("Error marshaling RepoDataArray: %v", err)
			return err
		}

		// Save the array in Badger using the given key
		err = txn.Set([]byte(repoKey), data)
		if err != nil {
			log.Printf("Error saving data to Badger for key %s: %v", repoKey, err)
			return err
		}

		log.Printf("RepoDataArray saved to Badger with key %s", repoKey)
		return nil
	})
}

// GetRepoArray retrieves an array of RepoData from Badger using the provided repoKey.
func (b *BadgerRepository) GetRepoArray(repoKey string) ([]domain.RepoData, error) {
	var repoDataArray []domain.RepoData

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(repoKey))
		if err != nil {
			log.Printf("Error retrieving data from Badger for key %s: %v", repoKey, err)
			return err
		}

		// Get the value and unmarshal it into the array of RepoData
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &repoDataArray)
		})
		if err != nil {
			log.Printf("Error unmarshaling data for key %s: %v", repoKey, err)
			return err
		}
		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			log.Printf("Repository data not found for key %s", repoKey)
			return nil, fmt.Errorf("repository data not found")
		}
		log.Printf("Error getting repo array: %v", err)
		return nil, err
	}

	log.Printf("RepoDataArray retrieved for key %s", repoKey)
	return repoDataArray, nil
}

// UpdateRepoArray updates or adds a RepoData entry in the array for the provided repoKey.
func (b *BadgerRepository) UpdateRepoArray(repoKey string, updatedRepo domain.RepoData) error {
	// Fetch the current array from Badger
	repoDataArray, err := b.GetRepoArray(repoKey)
	if err != nil {
		log.Printf("Error fetching repo array for update: %v", err)
		return err
	}

	// Add or update the RepoData in the array
	updated := false
	for i, repo := range repoDataArray {
		if repo.RepoName == updatedRepo.RepoName {
			// Update the existing repo entry
			repoDataArray[i] = updatedRepo
			updated = true
			break
		}
	}

	if !updated {
		// If not found, append the new repo entry
		repoDataArray = append(repoDataArray, updatedRepo)
	}

	// Save the updated array back to Badger
	err = b.SaveRepoArray(repoKey, repoDataArray)
	if err != nil {
		log.Printf("Error saving updated repo array: %v", err)
		return err
	}

	log.Printf("RepoDataArray updated successfully for key %s", repoKey)
	return nil
}

// Close closes the Badger database.
func (b *BadgerRepository) Close() error {
	err := b.db.Close()
	if err != nil {
		log.Printf("Error closing Badger database: %v", err)
		return err
	}
	log.Println("Badger database closed successfully")
	return nil
}
