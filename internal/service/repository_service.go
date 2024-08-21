package service

import (
	"github-service/internal/models"
	"github-service/internal/repository"
	"github-service/pkg/github"
	"gorm.io/gorm"
)

func RepositoryService(repo string, db *gorm.DB) (*models.Repository, error) {

	r, err := github.FetchRepositoryData(repo)
	if err != nil {

	}

	repository.SaveRepositories(db, r)

	return r, err
}
