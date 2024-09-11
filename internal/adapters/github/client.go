package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github-service/config"
	"github-service/pkg/httpClient"
	"github-service/pkg/logger"
)

// GithubClient provides methods to interact with the GitHub API
type GithubClient struct {
	client *httpclient.Client
	cfg    *config.Config
}

// NewGithubClient creates a new instance of GithubClient with a custom HTTP client.
// The client is rate-limited based on the configured PollInterval.
func NewGithubClient(cfg *config.Config, ctx context.Context) *GithubClient {
	rateLimitInterval := time.Duration(cfg.POLL_INTERVAL) * time.Second
	// Initialize the custom HTTP client
	client := httpclient.NewClient(nil, rateLimitInterval) // Use default http.Client or pass a custom one
	return &GithubClient{
		client: client,
		cfg:    cfg,
	}
}

// FetchRepositoryCommits fetches commits for a given repository from GitHub.
// It takes the repository owner, repository name, and the time from which to fetch commits.
func (g *GithubClient) FetchRepositoryCommits(ctx context.Context, owner, repo string, since time.Time) ([]Commit, error) {
	// Format the time as a string in RFC3339 format
	sinceStr := since.Format(time.RFC3339)
	// Construct the URL for fetching commits
	url := fmt.Sprintf("%s/%s/%s/commits?per_page=%s&until=%s", g.cfg.BASE_URL, owner, repo, g.cfg.PER_PAGE, sinceStr)

	// Perform the GET request using the custom HTTP client
	body, err := g.client.ApiCall(ctx, "GET", url, nil)
	if err != nil {
		logger.LogWarning(fmt.Sprintf("Error fetching commits for %s/%s: %v", owner, repo, err))
		return nil, err
	}

	// Unmarshal the response body into the slice of Commit structs
	var commits []Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		logger.LogWarning(fmt.Sprintf("Error unmarshaling commits for %s/%s: %v", owner, repo, err))
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("Fetched commits from %s/%s successfully", owner, repo))
	return commits, nil
}

// FetchRepositoryMetaData fetches metadata for a given repository from GitHub.
// It returns a Repository struct populated with metadata about the repository.
func (g *GithubClient) FetchRepositoryMetaData(ctx context.Context, owner, repo string) (*Repository, error) {
	// Construct the URL for fetching repository metadata
	url := fmt.Sprintf("%s/%s/%s", g.cfg.BASE_URL, owner, repo)

	// Perform the GET request using the custom HTTP client
	body, err := g.client.ApiCall(ctx, "GET", url, nil)
	if err != nil {
		logger.LogWarning(fmt.Sprintf("Error fetching metadata for repository %s/%s: %v", owner, repo, err))
		return nil, err
	}

	// Unmarshal the response body into the Repository struct
	var repository Repository
	if err := json.Unmarshal(body, &repository); err != nil {
		logger.LogWarning(fmt.Sprintf("Error unmarshaling repository metadata for %s/%s: %v", owner, repo, err))
		return nil, err
	}

	logger.LogInfo(fmt.Sprintf("Fetched repository metadata for %s/%s successfully", owner, repo))
	return &repository, nil
}
