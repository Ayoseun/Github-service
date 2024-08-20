package handlers

import (
	"github-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTopNCommitAuthors(c *gin.Context, db *gorm.DB) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil || n <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of authors"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	authors, err := repository.GetTopNCommitAuthors(db, n, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}
