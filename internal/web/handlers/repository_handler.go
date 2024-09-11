package handlers

import (
	"github-service/internal/core/domain"
	"github-service/internal/core/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RepositoryHandler handles HTTP requests related to repositories
type RepositoryHandler struct {
	repositoryService *service.RepositoryService
	monitorService    *service.MonitorService
}

// NewRepositoryHandler creates a new instance of RepositoryHandler with the given services
func NewRepositoryHandler(repositoryService *service.RepositoryService, monitorService *service.MonitorService) *RepositoryHandler {
	return &RepositoryHandler{
		repositoryService: repositoryService,
		monitorService:    monitorService,
	}
}

// FetchRepositoryData retrieves data for a given repository
func (h *RepositoryHandler) FetchRepositoryData(c *gin.Context) {
	// Get the repository name from the request parameters
	repo := c.Param("repo")

	// Call the RepositoryService to fetch the repository data
	repoData, err := h.repositoryService.GetRepository(c, repo)
	if err != nil {
		// Return 404 Not Found if the repository is not found
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": err.Error()})
		return
	}

	// Return the fetched repository data as JSON
	c.JSON(http.StatusOK, repoData)
}

// AddRepositoryToMonitor adds a repository to the commit monitor
func (h *RepositoryHandler) AddRepositoryToMonitor(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Query("repo")
	startDateStr := c.Query("start_date")

	// Parse the start date from the query string
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		// Return 400 Bad Request if the date format is invalid
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Invalid start date format"})
		return
	}

	// Create RepoData object
	repoData := domain.RepoData{
		Owner:    owner,
		RepoName: repoName,
	}

	// Add the repository to the monitor
	if err := h.monitorService.AddRepositoryCommitsToMonitor(c, repoData, startDate); err != nil {
		// Return 500 Internal Server Error if there is an issue adding the repository
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Failed to add repository"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository added successfully"})
}

// DeleteRepository removes a repository from the commit monitor
func (h *RepositoryHandler) DeleteRepository(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Query("repo")

	// Validate input parameters
	if owner == "" || repoName == "" {
		// Return 400 Bad Request if the owner or repository name is missing
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Owner and repository name are required"})
		return
	}

	// Remove the repository from the monitor
	ok, err := h.repositoryService.DeleteARepository(c, owner, repoName)
	if !ok {
		// Return 400 Bad Request if there's an issue removing the repository
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository removed successfully"})
}
