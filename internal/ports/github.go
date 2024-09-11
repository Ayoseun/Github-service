package ports

import (
	"context"
	"github-service/internal/core/domain"
	"time"
)

type GithubImpl interface {
	// FetchRepository fetches repository information for the given owner and repoName
	// Returns a Repository domain object and an error if the request fails
	FetchRepository(ctx context.Context, owner, repoName string) (*domain.Repository, error)

	// FetchCommit fetches a list of commits for the specified owner and repo starting from the 'since' time
	// Returns a slice of Commit domain objects and an error if the request fails
	FetchCommit(ctx context.Context, owner, repo string, since time.Time) ([]domain.Commit, error)
}
