package domain

import "time"

type Commit struct {
	Hash       string    `json:"sha"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	CommitDate time.Time `json:"date"`
	Email      string    `json:"email"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
}

// PaginatedResponse is the response structure for paginated commit data
type PaginatedResponse struct {
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
	Commits     []Commit `json:"commits"`
}
