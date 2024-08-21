package routes

import (
	"github-service/internal/handlers" // Importing the handlers package from the project's internal directory
	"github.com/gin-gonic/gin"         // Importing the Gin web framework
	"gorm.io/gorm"                     // Importing the GORM (Object-Relational Mapping) library for database interactions
)

// APPRoutes sets up the API routes for the application
func APPRoutes(r *gin.Engine, db *gorm.DB) {
	// Uncomment the following block if you want to set up a route to fetch repository commits
	// r.GET("/:repo/fetch_commits", func(c *gin.Context) {
	// 	handlers.FetchRepositoryCommits(c, db)
	// })

	// Set up a route to fetch repository data
	r.GET("/fetch_repository/:repo", func(c *gin.Context) {
		handlers.FetchRepositoryData(c, db)
	})

	// Set up a route to get the top N commit authors
	r.GET("/top_authors/:n", func(c *gin.Context) {
		handlers.GetTopNCommitAuthors(c, db)
	})

	// Set up a route to retrieve commits by repository
	r.GET("/commits/:repo", func(c *gin.Context) {
		handlers.RetrieveCommitsByRepository(c, db)
	})
}
