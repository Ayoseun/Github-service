package repository_test

import (
	"github-service/internal/models"
	"github-service/internal/web/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Automatically migrate the Repository and SavedCommit models to create the tables
	db.AutoMigrate(&models.Repository{}, &models.SavedCommit{})

	return db
}

func TestRetrieveCommitsByRepository(t *testing.T) {
	// Set up the in-memory SQLite database
	db := setupTestDB()

	// Seed the database with test data
	db.Create(&models.Repository{
		Name: "test-repo",
		URL:  "http://example.com/test-repo",
	})

	db.Create(&models.SavedCommit{
		Message: "Initial commit",
		Author:  "Author1",
		Date:    time.Now(),
		URL:     "http://example.com/test-repo/commit1",
	})

	tests := []struct {
		name           string
		repo           string
		page           string
		limit          string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Successful retrieval",
			repo:           "test-repo",
			page:           "1",
			limit:          "10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Repository not found",
			repo:           "non-existent-repo",
			page:           "1",
			limit:          "10",
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "repository non-existent-repo not found"},
		},
		{
			name:           "Invalid page number",
			repo:           "test-repo",
			page:           "-1",
			limit:          "10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid limit number",
			repo:           "test-repo",
			page:           "1",
			limit:          "-1",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.GET("/repositories/:repo/commits", func(c *gin.Context) {
				handlers.RetrieveCommitsByRepository(c, db)
			})

			req, _ := http.NewRequest("GET", "/repositories/"+tt.repo+"/commits?page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
