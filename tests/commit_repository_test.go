package repository_test

import (
	"github-service/internal/database/repository"
	"github-service/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

func cSetupTestDB() *gorm.DB {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Automatically migrate the SavedCommit model to create the table
	db.AutoMigrate(&models.SavedCommit{})

	return db
}
func TestSaveCommits(t *testing.T) {
	db := cSetupTestDB()
	repo := repository.NewCommitRepository(db)

	commit := &models.SavedCommit{
		Message: "Initial commit",
		Author:  "Author1",
		Date:    time.Now(),
		URL:     "http://example.com/commit1",
	}

	// Save the commit
	err := repo.SaveCommits(commit)
	assert.NoError(t, err)

	// Retrieve the commit from the database to verify it was saved correctly
	var savedCommit models.SavedCommit
	err = db.First(&savedCommit).Error
	assert.NoError(t, err)

	// Assert that the saved commit matches the original commit
	assert.Equal(t, commit.Message, savedCommit.Message)
	assert.Equal(t, commit.Author, savedCommit.Author)
	assert.Equal(t, commit.URL, savedCommit.URL)
}

func TestGetCommits(t *testing.T) {
	db := cSetupTestDB()
	repo := repository.NewCommitRepository(db)

	// Insert some commits into the database for testing
	commits := []models.SavedCommit{
		{Message: "Commit 1", Author: "Author1", Date: time.Now(), URL: "http://example.com/repo/commit1"},
		{Message: "Commit 2", Author: "Author2", Date: time.Now(), URL: "http://example.com/repo/commit2"},
	}
	for _, c := range commits {
		err := repo.SaveCommits(&c)
		assert.NoError(t, err)
	}

	// Test retrieving the commits with pagination
	retrievedCommits, err := repo.GetCommits("http://example.com/repo", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, retrievedCommits, 2)

	// Assert that the retrieved commits match the inserted commits
	assert.Equal(t, commits[0].Message, retrievedCommits[0].Message)
	assert.Equal(t, commits[1].Message, retrievedCommits[1].Message)
}

func TestGetTotalCommits(t *testing.T) {
	db := cSetupTestDB()
	repo := repository.NewCommitRepository(db)

	// Insert some commits into the database for testing
	commits := []models.SavedCommit{
		{Message: "Commit 1", Author: "Author1", Date: time.Now(), URL: "http://example.com/repo/commit1"},
		{Message: "Commit 2", Author: "Author2", Date: time.Now(), URL: "http://example.com/repo/commit2"},
	}
	for _, c := range commits {
		err := repo.SaveCommits(&c)
		assert.NoError(t, err)
	}

	// Test retrieving the total number of commits for the repository
	totalCommits, err := repo.GetTotalCommits("http://example.com/repo")
	assert.NoError(t, err)

	// Assert that the total number of commits matches the expected value
	assert.Equal(t, int64(2), totalCommits)
}

func TestGetCommitsPagination(t *testing.T) {
	db := cSetupTestDB()
	repo := repository.NewCommitRepository(db)

	// Insert more commits to test pagination
	for i := 1; i <= 15; i++ {
		commit := &models.SavedCommit{
			Message: "Commit " + strconv.Itoa(i),
			Author:  "Author" + strconv.Itoa(i),
			Date:    time.Now(),
			URL:     "http://example.com/repo/commit" + strconv.Itoa(i),
		}
		err := repo.SaveCommits(commit)
		assert.NoError(t, err)
	}

	// Test retrieving the first page with a limit of 10
	retrievedCommits, err := repo.GetCommits("http://example.com/repo", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, retrievedCommits, 10)

	// Test retrieving the second page with a limit of 10
	retrievedCommits, err = repo.GetCommits("http://example.com/repo", 2, 10)
	assert.NoError(t, err)
	assert.Len(t, retrievedCommits, 5)
}
