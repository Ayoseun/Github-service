package models

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	Language    string `json:"language"`
	ForksCount  int    `json:"forks_count"`
	StarsCount  int    `json:"stargazers_count"`
	OpenIssues  int    `json:"open_issues_count"`
	Watchers    int    `json:"watchers_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
