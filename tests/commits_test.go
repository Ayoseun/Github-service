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

func TestSaveCommit(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&models.SavedCommit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := repository.NewCommitRepository(db)
	assert.NoError(t, err)

	// Create a test commit
	testCommit := &models.SavedCommit{
		URL:        "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		Message:    "Fix all the bugs",
		Author:     "Monalisa Octocat",
		Date:       time.Date(2011, 4, 14, 16, 0, 49, 0, time.UTC),
		Repository: "Hello-World",
	}

	// Save the commit
	err = repo.SaveCommit(testCommit)
	assert.NoError(t, err)

	// Retrieve the saved commit
	var savedCommit models.SavedCommit
	result := db.First(&savedCommit)
	assert.NoError(t, result.Error)

	// Check if only one record was saved
	var count int64
	db.Model(&models.SavedCommit{}).Count(&count)
	assert.Equal(t, int64(1), count, "Expected only one record to be saved")

	// Validate the saved data
	assert.Equal(t, testCommit.URL, savedCommit.URL)
	assert.Equal(t, testCommit.Message, savedCommit.Message)
	assert.Equal(t, testCommit.Author, savedCommit.Author)
	assert.Equal(t, testCommit.Date.Unix(), savedCommit.Date.Unix())
	assert.Equal(t, testCommit.Repository, savedCommit.Repository)
}

func TestMultipleCommits(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&models.SavedCommit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := repository.NewCommitRepository(db)
	assert.NoError(t, err)

	// Create test commits
	testCommits := []models.SavedCommit{
		{
			URL:        "https://api.github.com/repos/octocat/Hello-World/commits/commit1",
			Message:    "First commit",
			Author:     "Alice",
			Date:       time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			Repository: "Hello-World",
		},
		{
			URL:        "https://api.github.com/repos/octocat/Hello-World/commits/commit2",
			Message:    "Second commit",
			Author:     "Bob",
			Date:       time.Date(2023, 1, 2, 11, 0, 0, 0, time.UTC),
			Repository: "Hello-World",
		},
		{
			URL:        "https://api.github.com/repos/octocat/Another-Repo/commits/commit3",
			Message:    "Third commit",
			Author:     "Charlie",
			Date:       time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
			Repository: "Another-Repo",
		},
	}

	// Save all commits
	for _, commit := range testCommits {
		err = repo.SaveCommit(&commit)
		assert.NoError(t, err)
	}

	// Test GetTotalCommits
	totalHelloWorld, err := repo.GetTotalCommits("Hello-World")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalHelloWorld)

	totalAnotherRepo, err := repo.GetTotalCommits("Another-Repo")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalAnotherRepo)

	// Test GetCommits (pagination)
	commitsHelloWorld, err := repo.GetCommits("Hello-World", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorld, 2)

	commitsAnotherRepo, err := repo.GetCommits("Another-Repo", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, commitsAnotherRepo, 1)

	// Validate the data of retrieved commits
	assert.Equal(t, testCommits[0].Message, commitsHelloWorld[0].Message)
	assert.Equal(t, testCommits[1].Message, commitsHelloWorld[1].Message)
	assert.Equal(t, testCommits[2].Message, commitsAnotherRepo[0].Message)

	// Test pagination
	commitsHelloWorldPage1, err := repo.GetCommits("Hello-World", 1, 1)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorldPage1, 1)
	assert.Equal(t, testCommits[0].Message, commitsHelloWorldPage1[0].Message)

	commitsHelloWorldPage2, err := repo.GetCommits("Hello-World", 2, 1)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorldPage2, 1)
	assert.Equal(t, testCommits[1].Message, commitsHelloWorldPage2[0].Message)
}
