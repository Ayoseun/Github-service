package handlers

import (
	"github-service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RepositoryHandler struct {
	repositoryService *service.RepositoryService
	commitMonitor     *service.CommitMonitor
}

func NewRepositoryHandler(repositoryService *service.RepositoryService, commitMonitor *service.CommitMonitor) *RepositoryHandler {
	return &RepositoryHandler{
		repositoryService: repositoryService,
		commitMonitor:     commitMonitor,
	}
}

// FetchRepositoryData is a Gin handler that fetches data for a given repository
func (h *RepositoryHandler) FetchRepositoryData(c *gin.Context) {
	// Get the repository name from the request parameters
	repo := c.Param("repo")
	// Call the RepositoryService to fetch the repository data
	r, err := h.repositoryService.GetRepository(repo)

	if err != nil {
		if err.Error() == "record not found" {
			// If the repository is not found, return a 404 Not Found error
			c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "Repository not found"})
		} else {
			// For other errors, return a 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Server error"})
		}
		return
	}
	// Return the fetched repository data
	c.JSON(http.StatusOK, r)
}

// AddRepositoryToMonitor adds a repository to the commit monitor
func (h *RepositoryHandler) AddRepositoryToMonitor(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Query("repo")
	startDateStr := c.Query("start_date")

	// Parse the start date
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Invalid start date format"})
		return
	}

	// Add the repository to the monitor
	if _, err := h.commitMonitor.AddRepository(owner, repo, startDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Failed to add repository"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository added successfully"})
}

// RemoveRepositoryFromMonitor removes a repository from the commit monitor
func (h *RepositoryHandler) RemoveRepositoryFromMonitor(c *gin.Context) {
	owner := c.Query("owner")
	name := c.Query("name")

	// Validate input
	if owner == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Owner and name are required"})
		return
	}

	// Remove the repository from the monitor
	h.commitMonitor.RemoveRepository(owner, name)

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository removed successfully"})
}

// FetchCommitsInRange fetches commits for a repository within a specified date range
func (h *RepositoryHandler) FetchCommitsInRange(c *gin.Context) {
	owner := c.Query("owner")
	repo := c.Query("repo")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Validate input
	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Owner and repository name are required"})
		return
	}

	// Parse the start date
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Invalid start date format"})
		return
	}

	var ok bool

	// Parse the end date if provided
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Invalid end date format"})
		return
	}
	// Call AddRepository with the date range
	ok, err = h.commitMonitor.AddRepository(owner, repo, startDate, endDate)

	if err != nil || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": "Failed to add repository for monitoring"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Repository monitoring started successfully"})
}
