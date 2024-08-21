package handlers

import (
	"github-service/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetTopNCommitAuthors is a Gin handler that retrieves the top N commit authors
func GetTopNCommitAuthors(c *gin.Context, db *gorm.DB) {
	rRepo := repository.NewRepository(db)
	// Parse the "n" parameter from the request
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		// If the "n" parameter is invalid, return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	// Parse the "page" and "limit" parameters from the query string
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Call the GetTopNCommitAuthors function in the repository package
	authors, err := rRepo.GetTopNCommitAuthors(n, page, limit)
	if err != nil {
		// If there's an error retrieving the top authors, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top authors"})
		return
	}

	// Return the top authors in the response
	c.JSON(http.StatusOK, authors)
}
