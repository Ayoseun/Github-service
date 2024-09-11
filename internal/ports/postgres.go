package ports

import (
	"context"
	"github-service/internal/core/domain"
)

// PostgresCommit defines the interface for commit data operations in a PostgreSQL database.
type PostgresCommit interface {
	// SaveCommit saves a commit to the database.
	// It returns an error if the save operation fails.
	SaveCommit(ctx context.Context, commit *domain.Commit) error

	// GetCommits retrieves a list of commits based on the repository URL, page number, and limit.
	// The page and limit parameters control pagination.
	// It returns a slice of commits and an error if the query fails.
	GetCommits(ctx context.Context, repositoryURL string, page, limit int) ([]domain.Commit, error)

	// GetTotalCommits retrieves the total number of commits for the specified repository name.
	// It returns the total count of commits and an error if the query fails.
	GetTotalCommits(ctx context.Context, repositoryName string) (int64, error)

	// DeleteAllCommits deletes all commits for the given repository name.
	// It returns a boolean indicating success and an error if the delete operation fails.
	DeleteAllCommits(ctx context.Context, repositoryName string) (bool, error)

	// GetLastCommitByRepositoryName retrieves the most recent commit for the specified repository based on the commit date.
	// It returns the latest commit and an error if the query fails or if no commits are found.
	GetLastCommitByRepositoryName(ctx context.Context, repoName string) (*domain.Commit, error)
}

// PostgresRepository defines the interface for repository data operations in a PostgreSQL database.
type PostgresRepository interface {
	// SaveRepository saves a repository to the database, creating a new entry if it doesn't exist or updating an existing one.
	// It returns an error if the save or update operation fails.
	SaveRepository(ctx context.Context, repo *domain.Repository) error

	// GetTopNCommitAuthors retrieves the top N commit authors for the specified repository, with pagination support.
	// It groups authors by name, counts their commits, and orders the results by the count in descending order.
	// It returns a slice of top authors with their commit counts and an error if the query fails.
	GetTopNCommitAuthors(ctx context.Context, repository string, page, limit int) (domain.TopAuthorsCount, error)

	// GetRepositoryByName retrieves a repository based on its name.
	// It returns the repository model and an error if the query fails or if the repository is not found.
	GetRepositoryByName(ctx context.Context, repository string) (domain.Repository, error)

	// DeleteRepository deletes a repository with the given name and owner.
	// It returns a boolean indicating whether the deletion was successful and an error if the delete operation fails.
	DeleteRepository(ctx context.Context, owner, repositoryName string) (bool, error)
}
