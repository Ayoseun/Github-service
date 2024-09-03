package repository_test

import (
	"github-service/internal/database/repository"
	"github-service/internal/domain/models"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTopNCommitAuthors(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&models.SavedCommit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := repository.NewRepository(db)
	assert.NoError(t, err)

	// Create test commits
	testCommits := []models.SavedCommit{
		{Author: "Alice", Date: time.Now()},
		{Author: "Alice", Date: time.Now()},
		{Author: "Bob", Date: time.Now()},
		{Author: "Charlie", Date: time.Now()},
		{Author: "Alice", Date: time.Now()},
		{Author: "Bob", Date: time.Now()},
		{Author: "David", Date: time.Now()},
	}

	// Save all commits
	for _, commit := range testCommits {
		err = db.Create(&commit).Error
		assert.NoError(t, err)
	}

	// Test cases
	testCases := []struct {
		name          string
		page          int
		limit         int
		expectedCount int
		expectedOrder []string
	}{
		{
			name:          "First page, limit 2",
			page:          1,
			limit:         2,
			expectedCount: 2,
			expectedOrder: []string{"Alice", "Bob"},
		},
		{
			name:          "Second page, limit 2",
			page:          2,
			limit:         2,
			expectedCount: 2,
			expectedOrder: []string{"Charlie", "David"},
		},
		{
			name:          "All authors, single page",
			page:          1,
			limit:         10,
			expectedCount: 4,
			expectedOrder: []string{"Alice", "Bob", "Charlie", "David"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authors, err := repo.GetTopNCommitAuthors(tc.page, tc.limit)
			assert.NoError(t, err)
			assert.Len(t, authors, tc.expectedCount)

			for i, expectedAuthor := range tc.expectedOrder {
				if i < len(authors) {
					assert.Equal(t, expectedAuthor, authors[i].Author)
				}
			}

			// Check if authors are ordered by commit count
			for i := 1; i < len(authors); i++ {
				assert.GreaterOrEqual(t, authors[i-1].Count, authors[i].Count)
			}
		})
	}
}
