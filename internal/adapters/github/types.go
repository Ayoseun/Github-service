package github

import "time"

type Commit struct {
	Author      interface{} `json:"author"`
	CommentsURL string      `json:"comments_url"`
	Commit      struct {
		Author struct {
			Date  time.Time `json:"date"`
			Email string    `json:"email"`
			Name  string    `json:"name"`
		} `json:"author"`
		CommentCount int `json:"comment_count"`
		Committer    struct {
			Date  time.Time `json:"date"`
			Email string    `json:"email"`
			Name  string    `json:"name"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			SHA string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		Verification struct {
			Payload   interface{} `json:"payload"`
			Reason    string      `json:"reason"`
			Signature interface{} `json:"signature"`
			Verified  bool        `json:"verified"`
		} `json:"verification"`
	} `json:"commit"`
	Committer interface{} `json:"committer"`
	HTMLURL   string      `json:"html_url"`
	NodeID    string      `json:"node_id"`
	Parents   []struct {
		HTMLURL string `json:"html_url"`
		SHA     string `json:"sha"`
		URL     string `json:"url"`
	} `json:"parents"`
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Repository struct {
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
