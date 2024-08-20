package handlers

import (
	"github-service/internal/models"
	"github-service/internal/repository"
	"net/http"
	"strconv"

	"github-service/pkg/github"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchRepositoryCommits(c *gin.Context, db *gorm.DB) {
	repo := c.Param("repo")
	commits, err := github.FetchRepositoryCommits(repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch external data"})
		return
	}

	for _, commit := range commits {
		savedCommit := &models.SavedCommit{
			Message: commit.Commit.Message,
			Author:  commit.Commit.Author.Name,
			Date:    commit.Commit.Author.Date,
			URL:     commit.HTMLURL,
		}
		repository.SaveCommit(db, savedCommit)
	}

	c.JSON(http.StatusOK, commits)
}
func RetrieveCommitsByRepository(c *gin.Context, db *gorm.DB) {
	repo := c.Param("repo")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	var r models.Repository
	if err := db.Where("name = ?", repo).First(&r).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	commits, err := repository.GetCommitsByRepository(db, r.URL, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commits"})
		return
	}

	c.JSON(http.StatusOK, commits)
}
