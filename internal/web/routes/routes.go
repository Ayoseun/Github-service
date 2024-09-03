package routes

import (
	"github-service/internal/web/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes sets up the API routes for the application
func SetupAPIRoutes(r *gin.Engine, commitHandler *handlers.CommitHandler, repositoryHandler *handlers.RepositoryHandler) {
	// Route to fetch repository data
	r.GET("/repositories/:repo/fetch", repositoryHandler.FetchRepositoryData)

	// Route to get the top N commit authors
	r.GET("/repositories/:repo/top-authors/:n", commitHandler.GetTopNCommitAuthors)

	// Route to retrieve commits by repository
	r.GET("/repositories/:repo/commits", commitHandler.GetCommits)
}
