package service

import (
	"github-service/internal/models"
	"github-service/internal/repository"
	"github-service/pkg/github"
	"time"

	"gorm.io/gorm"
)

func CommitsService(repo string, db *gorm.DB, lastFetchedCommitDate time.Time) []models.Commit {

	commits, err := github.FetchRepositoryCommits(repo)
	if err != nil {

	}

	for _, commit := range commits {
		if commit.Commit.Author.Date.After(lastFetchedCommitDate) {
			savedCommit := &models.SavedCommit{
				Message: commit.Commit.Message,
				Author:  commit.Commit.Author.Name,
				Date:    commit.Commit.Author.Date,
				URL:     commit.HTMLURL,
			}
			repository.SaveCommits(db, savedCommit)
		}
	}
	return commits
}
