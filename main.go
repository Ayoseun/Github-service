package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Commit struct {
	Author      interface{} `json:"author"`
	CommentsURL string      `json:"comments_url"`
	Commit      struct {
		Author struct {
			Date  string `json:"date"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"author"`
		CommentCount int `json:"comment_count"`
		Committer    struct {
			Date  string `json:"date"`
			Email string `json:"email"`
			Name  string `json:"name"`
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

type SavedCommit struct {
	gorm.Model
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

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

var db *gorm.DB

func initDB() {
	dsn := "postgres://ayoseun:Jared15$@localhost:5432/github"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&SavedCommit{}, &Repository{})
}

func fetchRepositoryCommits(c *gin.Context) {
	repo := c.Param("repo")
	url := fmt.Sprintf("https://api.github.com/repos/chromium/%s/commits", repo)
	println(url)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch external data"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch external data"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var commits []Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
		return
	}

	for _, commit := range commits {
		savedCommit := SavedCommit{
			Message: commit.Commit.Message,
			Author:  commit.Commit.Author.Name,
			Date:    commit.Commit.Author.Date,
			URL:     commit.HTMLURL,
		}
		db.Create(&savedCommit)
	}

	c.JSON(http.StatusOK, commits)
}

func fetchRepositoryData(c *gin.Context) {
	repo := c.Param("repo")
	url := fmt.Sprintf("https://api.github.com/repos/chromium/%s", repo)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository data"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository data"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var r Repository
	if err := json.Unmarshal(body, &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body"})
		return
	}

	db.Create(&r)

	c.JSON(http.StatusOK, r)
}

func getTopNCommitAuthors(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var authors []struct {
		Author string `json:"author"`
		Count  int    `json:"count"`
	}

	db.Model(&SavedCommit{}).
		Select("author, count(author) as count").
		Group("author").
		Order("count desc").
		Limit(limit).
		Offset(offset).
		Scan(&authors)

	c.JSON(http.StatusOK, authors)
}

func retrieveCommitsByRepository(c *gin.Context) {
	repo := c.Param("repo")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var repository Repository
	if err := db.Where("name = ?", repo).First(&repository).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var commits []SavedCommit
	if err := db.Where("url LIKE ?", fmt.Sprintf("%%%s%%", repository.URL)).
		Limit(limit).
		Offset(offset).
		Find(&commits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commits"})
		return
	}

	c.JSON(http.StatusOK, commits)
}

func main() {
	initDB()

	router := gin.Default()
	router.GET("/:repo/fetch_commits", fetchRepositoryCommits)
	router.GET("/fetch_repository/:repo", fetchRepositoryData)
	router.GET("/top_authors/:n", getTopNCommitAuthors)
	router.GET("/commits/:repo", retrieveCommitsByRepository)

	router.Run(":8080")
}
