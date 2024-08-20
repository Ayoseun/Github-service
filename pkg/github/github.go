package github

import (
	"encoding/json"
	"fmt"
	"github-service/internal/models"
	"io"
	"net/http"
)

func FetchRepositoryCommits(repo string) ([]models.Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/chromium/%s/commits", repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch external data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var commits []models.Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func FetchRepositoryData(repo string) (*models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/chromium/%s", repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r models.Repository
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
