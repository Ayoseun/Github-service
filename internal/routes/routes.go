package routes

import (
	"github-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func APPRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/:repo/fetch_commits", func(c *gin.Context) {
		handlers.FetchRepositoryCommits(c, db)
	})

	r.GET("/fetch_repository/:repo", func(c *gin.Context) {
		handlers.FetchRepositoryData(c, db)
	})

	r.GET("/top_authors/:n", func(c *gin.Context) {
		handlers.GetTopNCommitAuthors(c, db)
	})

	r.GET("/commits/:repo", func(c *gin.Context) {
		handlers.RetrieveCommitsByRepository(c, db)
	})
}
