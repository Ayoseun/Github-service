package handlers

import (
	"errors"
	"fmt"
	"github-service/internal/models"
	"github-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrRepositoryNotFound is a custom error type used when a repository is not found
type ErrRepositoryNotFound struct {
	Repo string
}

func (e ErrRepositoryNotFound) Error() string {
	return fmt.Sprintf("repository %s not found", e.Repo)
}

// RetrieveCommitsByRepository is a Gin handler that retrieves commits for a given repository
func RetrieveCommitsByRepository(c *gin.Context, db *gorm.DB) {
	// Create a new instance of the CommitRepository
	commitRepo := repository.NewCommitRepository(db)

	// Get the repository name from the request parameters
	repo := c.Param("repo")

	// Parse the page and limit parameters from the query string
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Retrieve the repository from the database
	var r models.Repository
	if err := db.Where("name = ?", repo).First(&r).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the repository is not found, return a 404 error with a custom error message
			c.JSON(http.StatusNotFound, gin.H{"error": ErrRepositoryNotFound{repo}.Error()})
			return
		}
		// If there's another error, return a 500 error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve repository"})
		return
	}

	// Retrieve the total number of commits for the repository
	totalCommits, err := commitRepo.GetTotalCommits(r.URL)
	if err != nil {
		// If there's an error retrieving the total commits, return a 500 error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve total commits"})
		return
	}

	// Calculate the total number of pages based on the total commits and the requested limit
	totalPages := int((totalCommits + int64(limit) - 1) / int64(limit))

	// Retrieve the commits for the given page and limit
	commits, err := commitRepo.GetCommits(r.URL, page, limit)
	if err != nil {
		// If there's an error retrieving the commits, return a 500 error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve commits"})
		return
	}

	// Return the paginated commit data
	c.JSON(http.StatusOK, models.PaginatedResponse{
		CurrentPage: page,
		TotalPages:  totalPages,
		Commits:     commits,
	})
}
