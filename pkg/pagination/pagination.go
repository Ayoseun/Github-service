package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParsePaginationParams parses pagination parameters from the query string
func ParsePaginationParams(c *gin.Context) (page int, limit int, err error) {
	page, err = strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	return page, limit, nil
}

// RespondWithError handles error responses
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"statusCode": statusCode, "message": message})
}
