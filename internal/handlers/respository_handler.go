package handlers

import (
	"github-service/internal/repository"
	"github-service/pkg/github"
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func FetchRepositoryData(c *gin.Context, db *gorm.DB) {
	repo := c.Param("repo")
	r, err := github.FetchRepositoryData(repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository data"})
		return
	}

	repository.SaveRepository(db, r)

	c.JSON(http.StatusOK, r)
}
