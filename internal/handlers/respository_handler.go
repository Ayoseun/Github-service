package handlers

import (
	"github-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// FetchRepositoryData is a Gin handler that fetches data for a given repository
func FetchRepositoryData(c *gin.Context, db *gorm.DB) {
	// Get the repository name from the request parameters
	repo := c.Param("repo")

	// Call the RepositoryService to fetch the repository data
	r, err := service.RepositoryService(repo, db)
	if err != nil {
		// If there's an error fetching the repository data, return a 500 error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository data"})
		return
	}

	// Return the fetched repository data
	c.JSON(http.StatusOK, r)
}
