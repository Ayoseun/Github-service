package domain

import "time"

type Repository struct {
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
type TopAuthorsCount []struct {
	Author string `json:"author"`
	Count  int    `json:"count"`
}

type RepoData struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repoName"`
}
