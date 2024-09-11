package postgresdb

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	Owner            string    `json:"owner"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	URL              string    `json:"html_url"`
	Language         string    `json:"language"`
	ForksCount       int       `json:"forks_count"`
	StarsGazersCount int       `json:"stargazers_count"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	WatchersCount    int       `json:"watchers_count"`
	SubscribersCount int       `json:"subscribers_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Commit struct {
	gorm.Model
	Hash       string    `json:"sha"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	CommitDate time.Time `json:"date"`
	Email      string    `json:"email"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
}
