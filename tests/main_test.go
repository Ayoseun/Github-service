package test_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTopNCommitAuthors(t *testing.T) {
	// Set up the test database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&SavedCommit{})

	// Insert test data
	commits := []SavedCommit{
		{Message: "Commit 1", Author: "Author A", Date: "2023-01-01", URL: "http://example.com/commit1"},
		{Message: "Commit 2", Author: "Author B", Date: "2023-01-02", URL: "http://example.com/commit2"},
		{Message: "Commit 3", Author: "Author A", Date: "2023-01-03", URL: "http://example.com/commit3"},
		{Message: "Commit 4", Author: "Author C", Date: "2023-01-04", URL: "http://example.com/commit4"},
	}
	db.Create(&commits)

	// Set up the Gin router
	router := gin.Default()
	router.GET("/top_authors/:n", func(c *gin.Context) {
		getTopNCommitAuthors(c)
	})

	// Create a request to the endpoint
	req, _ := http.NewRequest("GET", "/top_authors/2?page=1&limit=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `[{"author":"Author A","count":2},{"author":"Author B","count":1}]`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}


