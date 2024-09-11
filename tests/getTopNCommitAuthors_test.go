package repository_test

import (
	"context"
	"github-service/internal/adapters/postgresdb"
	"github-service/internal/core/domain"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTopNCommitAuthors(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	ctx := context.Background()
	// Auto migrate the schema
	err = db.AutoMigrate(&postgresdb.Commit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := postgresdb.NewRepository(db)
	assert.NoError(t, err)

	// Create test commits
	testCommits := []domain.Commit{
		{Author: "Alice", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "Alice", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "Bob", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "Charlie", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "Alice", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "Bob", CommitDate: time.Now(), Repository: "Hello-World"},
		{Author: "David", CommitDate: time.Now(), Repository: "Hello-World"},
	}

	// Save all commits
	for _, commit := range testCommits {
		err = db.Create(&commit).Error
		assert.NoError(t, err)
	}

	// Test cases
	testCases := []struct {
		repo          string
		name          string
		page          int
		limit         int
		expectedCount int
		expectedOrder []string
	}{
		{
			repo:          "Hello-World",
			name:          "First page, limit 2",
			page:          1,
			limit:         2,
			expectedCount: 2,
			expectedOrder: []string{"Alice", "Bob"},
		},
		{
			repo:          "Hello-World",
			name:          "Second page, limit 2",
			page:          2,
			limit:         2,
			expectedCount: 2,
			expectedOrder: []string{"Charlie", "David"},
		},
		{
			repo:          "Hello-World",
			name:          "All authors, single page",
			page:          1,
			limit:         10,
			expectedCount: 4,
			expectedOrder: []string{"Alice", "Bob", "Charlie", "David"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authors, err := repo.GetTopNCommitAuthors(ctx, tc.repo, tc.page, tc.limit)
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
