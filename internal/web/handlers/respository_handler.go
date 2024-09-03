package handlers

import (
	"github-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RepositoryHandler struct {
	repositoryService *service.RepositoryService
}

func NewRepositoryHandler(repositoryService *service.RepositoryService) *RepositoryHandler {
	return &RepositoryHandler{

		repositoryService: repositoryService,
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

// AddRepository is a Gin handler that adds a given repository Metadata
func (h *RepositoryHandler) AddRepository(c *gin.Context) {
	// Get the repository name from the request parameters
	repo := c.Query("repo")
	owner := c.Query("owner")
	// Call the RepositoryService to fetch the repository data
	r, err := h.repositoryService.FetchAndSaveRepository(owner, repo)

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
