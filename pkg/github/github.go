package github

import (
	"encoding/json"
	"fmt"
	"github-service/internal/config"
	"github-service/internal/domain/models"
	"io"
	"net/http"
)

func FetchRepositoryCommits(owner, repo string, cfg config.Config) ([]models.Commit, error) {
	url := fmt.Sprintf("%s/%s/%s/commits", cfg.BASE_URL, owner, repo)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the API key to the Authorization header only if it is provided
	if cfg.GITHUB_TOKEN != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", cfg.GITHUB_TOKEN))
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch external data: %s", resp.Status)
	}

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var commits []models.Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}

	return commits, err
}

func FetchRepositoryMetaData(owner, repo string, cfg config.Config) (*models.Repository, error) {
	url := fmt.Sprintf("%s/%s/%s", cfg.BASE_URL, owner, repo)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the API key to the Authorization header only if it is provided
	if cfg.GITHUB_TOKEN != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", cfg.GITHUB_TOKEN))
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch external data: %s", resp.Status)
	}

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r models.Repository

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return &r, err
}
