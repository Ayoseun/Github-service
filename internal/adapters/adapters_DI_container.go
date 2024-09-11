package adapters

import (
	"fmt"
	"github-service/config"
	badger "github-service/internal/adapters/badgerdb"
	"github-service/internal/adapters/postgresdb"
	"github-service/internal/core/domain"
	"github-service/internal/ports"
)

func SetupStorage(cfg config.Config) (ports.PostgresCommit, ports.PostgresRepository, *badger.BadgerRepository, error) {
	// Initialize the Postgres database connection
	db, err := postgresdb.Connect(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create the Commit repository
	commitRepo, err := postgresdb.NewCommitRepository(db)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create commit repository: %w", err)
	}

	// Create the Repository repository
	repositoryRepo, err := postgresdb.NewRepository(db)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create repository repository: %w", err)
	}

	// Initialize Badger key-value store
	badgerService, err := badger.NewBadgerRepository("./tmp")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create badger repository: %w", err)
	}

	// Example of seeding data into Badger (can be optional)
	repoDataArray := []domain.RepoData{
		{RepoName: "chromium", Owner: "chromium"},
	}
	err = badgerService.SaveRepoArray("repos", repoDataArray)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to save initial data to badger: %w", err)
	}

	return commitRepo, repositoryRepo, badgerService, nil
}
