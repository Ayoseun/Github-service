package handlers

import (
	"github-service/internal/core/domain"
	"github-service/internal/core/service"
	"github-service/pkg/pagination"
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
	repo := c.Param("repo")

	// Parse pagination parameters from the query string
	page, limit, err := pagination.ParsePaginationParams(c)
	if err != nil {
		pagination.RespondWithError(c, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	// Retrieve the repository details
	url, err := h.repositoryService.GetRepository(c, repo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": err.Error()})
		return
	}

	// Retrieve the total number of commits for the repository
	totalCommits, err := h.commitService.GetCommitCount(c, url.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total commits"})
		return
	}

	// Calculate total pages based on the total commits and limit
	totalPages := int((totalCommits + int64(limit) - 1) / int64(limit))
	if totalPages <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No commits found"})
		return
	}

	// Retrieve the commits for the requested page and limit
	commits, err := h.commitService.GetPaginatedCommits(c, url.Name, page, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No commits found"})
		return
	}

	// Return the paginated commit data as JSON
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data": domain.PaginatedResponse{
			CurrentPage: page,
			TotalPages:  totalPages,
			Commits:     commits,
		},
	})
}

// GetTopNCommitAuthors retrieves the top N commit authors and returns them as JSON
func (h *CommitHandler) GetTopNCommitAuthors(c *gin.Context) {
	repo := c.Param("repo")

	// Parse the "n" parameter to determine the number of top authors
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	page, limit, err := pagination.ParsePaginationParams(c)
	if err != nil {
		pagination.RespondWithError(c, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	// Retrieve the top N commit authors from the repository service
	authors, err := h.repositoryService.GetTopNCommitAuthors(c, repo, n, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		return
	}
	if len(authors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "No authors found"})
		return
	}

	// Return the top authors as JSON
	c.JSON(http.StatusOK, authors)
}

// ResetCollection removes all commits for a specific repository and returns a success message
func (h *CommitHandler) ResetCollection(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Query("repo")

	// Validate input
	if owner == "" || repoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Owner and repository name are required"})
		return
	}

	// Remove all commits for the specified repository
	ok, err := h.commitService.DeleteCommits(c, repoName)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository commits removed successfully"})
}
