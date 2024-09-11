package routes

import (
	"github-service/internal/web/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes sets up the API routes for the application.
func SetupAPIRoutes(r *gin.Engine, commitHandler *handlers.CommitHandler, repositoryHandler *handlers.RepositoryHandler) {

	// Route to fetch repository data
	// GET /repositories/:repo/fetch
	// Retrieves detailed information about a specific repository.
	r.GET("/repositories/:repo/fetch", repositoryHandler.FetchRepositoryData)

	// Route to get the top N commit authors
	// GET /repositories/:repo/top-authors/:n
	// Retrieves the top N commit authors for a specific repository.
	r.GET("/repositories/:repo/top-authors/:n", commitHandler.GetTopNCommitAuthors)

	// Route to retrieve commits for a repository
	// GET /repositories/:repo/commits
	// Retrieves a list of commits for a specific repository.
	r.GET("/repositories/:repo/commits", commitHandler.GetCommits)

	// Route to reset commits for a repository
	// GET /repositories/:repo/reset
	// Resets or clears commit data for a specific repository.
	r.GET("/repositories/:repo/reset", commitHandler.ResetCollection) // Renamed handler to reflect reset operation

	// Route to add a repository to the monitoring service
	// GET /repositories/monitor/:owner
	// Adds a repository to the monitoring service, including commit pulling history.
	r.GET("/repositories/monitor/:owner", repositoryHandler.AddRepositoryToMonitor)

	// Route to remove a repository from the monitoring service
	// DELETE /repositories/monitor/:owner
	// Removes a repository from the monitoring service.
	r.DELETE("/repositories/monitor/:owner", repositoryHandler.DeleteRepository)
}
