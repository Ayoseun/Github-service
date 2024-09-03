package repository_test

import (
	"github-service/internal/database/repository"
	"github-service/internal/domain/models"
	"testing"

	"sort"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Helper function to sort TopAuthorsCount
func sortTopAuthorsCount(authors models.TopAuthorsCount) {
	sort.Slice(authors, func(i, j int) bool {
		if authors[i].Count == authors[j].Count {
			return authors[i].Author < authors[j].Author
		}
		return authors[i].Count > authors[j].Count
	})
}

func TestGetTopNCommitAuthors(t *testing.T) {
	// Initialize an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	// Auto-migrate the SavedCommit model to create the necessary schema
	err = db.AutoMigrate(&models.SavedCommit{})
	assert.NoError(t, err)

	// Seed the database with test data
	commits := []models.SavedCommit{
		{Author: "Author1"},
		{Author: "Author2"},
		{Author: "Author1"},
		{Author: "Author3"},
		{Author: "Author1"},
		{Author: "Author2"},
		{Author: "Author3"},
		{Author: "Author4"},
		{Author: "Author5"},
		{Author: "Author5"},
		{Author: "Author6"},
	}

	for _, commit := range commits {
		db.Create(&commit)
	}

	// Create the repository
	repo, err := repository.NewRepository(db)
	assert.NoError(t, err)

	// Define the test cases
	testCases := []struct {
		page           int
		limit          int
		expectedResult models.TopAuthorsCount
		expectedLength int
	}{
		{
			page:  1,
			limit: 3,
			expectedResult: models.TopAuthorsCount{
				{Author: "Author1", Count: 3},
				{Author: "Author2", Count: 2},
				{Author: "Author3", Count: 2},
			},
			expectedLength: 3,
		},
		{
			page:  2,
			limit: 3,
			expectedResult: models.TopAuthorsCount{
				{Author: "Author4", Count: 1},
				{Author: "Author5", Count: 2},
				{Author: "Author6", Count: 1},
			},
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		// Call the method under test
		authors, err := repo.GetTopNCommitAuthors(tc.page, tc.limit)
		assert.NoError(t, err)

		// Sort the results
		sortTopAuthorsCount(authors)
		sortTopAuthorsCount(tc.expectedResult)

		// Assert the results match the expected data and length
		assert.Equal(t, tc.expectedResult, authors)
		assert.Len(t, authors, tc.expectedLength)
	}
}
