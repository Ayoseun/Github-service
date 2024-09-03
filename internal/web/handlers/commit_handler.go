package handlers

import (
	"fmt"
	"github-service/internal/domain/models"
	"github-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommitHandler handles HTTP requests related to commits
type CommitHandler struct {
	commitService     *service.CommitService
	repositoryService *service.RepositoryService
}

// NewCommitHandler creates a new instance of CommitHandler with the given services
func NewCommitHandler(commitService *service.CommitService, repositoryService *service.RepositoryService) *CommitHandler {
	return &CommitHandler{
		commitService:     commitService,
		repositoryService: repositoryService,
	}
}

// GetCommits retrieves commits for a given repository and returns them as a paginated response
func (h *CommitHandler) GetCommits(c *gin.Context) {
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

	// Retrieve the repository details from the database
	url, err := h.repositoryService.GetRepository(repo)
	if err != nil {
		if err.Error() == "record not found" {
			// Return 404 Not Found if the repository is not found
			c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "Repository not found"})
		} else {
			// Return 500 Internal Server Error for other errors
			c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		}
		return
	}

	// Retrieve the total number of commits for the repository
	totalCommits, err := h.commitService.GetCommitCount(url.Name)
	fmt.Println(totalCommits)
	if err != nil {
		// Return 500 Internal Server Error if there's an error retrieving the commit count
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve total commits"})
		return
	}

	// Calculate the total number of pages based on the total commits and limit
	totalPages := int((totalCommits + int64(limit) - 1) / int64(limit))

	// Retrieve the commits for the requested page and limit
	commits, err := h.commitService.GetPaginatedCommits(url.Name, page, limit)
	if err != nil {
		if len(commits) == 0 {
			// Return 404 Not Found if no commits are found
			c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No commits found"})
		} else {
			// Return 500 Internal Server Error for other errors
			c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		}
		return
	}

	// Return the paginated commit data as JSON
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusFound, "data": models.PaginatedResponse{
		CurrentPage: page,
		TotalPages:  totalPages,
		Commits:     commits,
	}})
}

// GetTopNCommitAuthors retrieves the top N commit authors and returns them as JSON
func (h *CommitHandler) GetTopNCommitAuthors(c *gin.Context) {
	// Parse the "n" parameter from the request to determine the number of top authors
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		// Return 400 Bad Request if the "n" parameter is invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	// Parse the page and limit parameters from the query string
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Retrieve the top N commit authors from the repository service
	authors, err := h.repositoryService.GetTopNCommitAuthors(n, page, limit)
	if err != nil {

		// Return 500 Internal Server Error for other errors
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})

		return
	}
	if len(authors) == 0 {
		// Return 404 Not Found if no authors are found
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No authors found"})
		return
	}
	// Return the top authors as JSON
	c.JSON(http.StatusOK, authors)
}
