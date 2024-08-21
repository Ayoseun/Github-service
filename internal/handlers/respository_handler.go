package handlers

import (
	"github-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func FetchRepositoryData(c *gin.Context, db *gorm.DB) {
	repo := c.Param("repo")
	r, err := service.RepositoryService(repo, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository data"})
		return
	}

	c.JSON(http.StatusOK, r)
}
