package repository_test

import (
	"context"
	"github-service/internal/adapters/postgresdb"
	"github-service/internal/core/domain"
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
	ctx := context.Background()
	// Auto migrate the schema
	err = db.AutoMigrate(&postgresdb.Commit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := postgresdb.NewCommitRepository(db)
	assert.NoError(t, err)

	// Create a test commit
	testCommit := &domain.Commit{
		URL:        "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		Message:    "Fix all the bugs",
		Author:     "Monalisa Octocat",
		CommitDate: time.Now(),
		Repository: "Hello-World",
	}

	// Save the commit
	err = repo.SaveCommit(ctx, testCommit)
	assert.NoError(t, err)

	// Retrieve the saved commit
	var savedCommit domain.Commit
	result := db.First(&savedCommit)
	assert.NoError(t, result.Error)

	// Check if only one record was saved
	var count int64
	db.Model(&domain.Commit{}).Count(&count)
	assert.Equal(t, int64(1), count, "Expected only one record to be saved")

	// Validate the saved data
	assert.Equal(t, testCommit.URL, savedCommit.URL)
	assert.Equal(t, testCommit.Message, savedCommit.Message)
	assert.Equal(t, testCommit.Author, savedCommit.Author)
	assert.WithinDuration(t, testCommit.CommitDate, savedCommit.CommitDate, time.Second, "Commit dates should be within 1 second of each other")
	assert.Equal(t, testCommit.Repository, savedCommit.Repository)
}

func TestMultipleCommits(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	ctx := context.Background()
	// Auto migrate the schema
	err = db.AutoMigrate(&postgresdb.Commit{})
	assert.NoError(t, err)

	// Create repository instance
	repo, err := postgresdb.NewCommitRepository(db)
	assert.NoError(t, err)

	// Create test commits
	testCommits := []domain.Commit{
		{
			URL:        "https://api.github.com/repos/octocat/Hello-World/commits/commit1",
			Message:    "First commit",
			Author:     "Alice",
			CommitDate: time.Now(),
			Repository: "Hello-World",
		},
		{
			URL:        "https://api.github.com/repos/octocat/Hello-World/commits/commit2",
			Message:    "Second commit",
			Author:     "Bob",
			CommitDate: time.Now(),
			Repository: "Hello-World",
		},
		{
			URL:        "https://api.github.com/repos/octocat/Another-Repo/commits/commit3",
			Message:    "Third commit",
			Author:     "Charlie",
			CommitDate: time.Now(),
			Repository: "Another-Repo",
		},
	}

	// Save all commits
	for _, commit := range testCommits {
		err = repo.SaveCommit(ctx, &commit)
		assert.NoError(t, err)
	}

	// Test GetTotalCommits
	totalHelloWorld, err := repo.GetTotalCommits(ctx, "Hello-World")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalHelloWorld)

	totalAnotherRepo, err := repo.GetTotalCommits(ctx, "Another-Repo")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalAnotherRepo)

	// Test GetCommits (pagination)
	commitsHelloWorld, err := repo.GetCommits(ctx, "Hello-World", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorld, 2)

	commitsAnotherRepo, err := repo.GetCommits(ctx, "Another-Repo", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, commitsAnotherRepo, 1)

	// Validate the data of retrieved commits
	assert.Equal(t, testCommits[0].Message, commitsHelloWorld[0].Message)
	assert.Equal(t, testCommits[1].Message, commitsHelloWorld[1].Message)
	assert.Equal(t, testCommits[2].Message, commitsAnotherRepo[0].Message)

	// Test pagination
	commitsHelloWorldPage1, err := repo.GetCommits(ctx, "Hello-World", 1, 1)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorldPage1, 1)
	assert.Equal(t, testCommits[0].Message, commitsHelloWorldPage1[0].Message)

	commitsHelloWorldPage2, err := repo.GetCommits(ctx, "Hello-World", 2, 1)
	assert.NoError(t, err)
	assert.Len(t, commitsHelloWorldPage2, 1)
	assert.Equal(t, testCommits[1].Message, commitsHelloWorldPage2[0].Message)
}
