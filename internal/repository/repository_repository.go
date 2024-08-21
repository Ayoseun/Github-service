package repository

import (
	"github-service/internal/models"

	"gorm.io/gorm"
)

func SaveRepositories(db *gorm.DB, repository *models.Repository) error {
	// Check if the repository already exists
	var existingRepo models.Repository
	result := db.Where("name = ?", repository.Name).First(&existingRepo)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Repository doesn't exist, create a new one
			return db.Create(repository).Error
		} else {
			// Other error occurred, return it
			return result.Error
		}
	}

	// Repository exists, update it
	existingRepo.Description = repository.Description
	existingRepo.URL = repository.URL
	existingRepo.Language = repository.Language
	existingRepo.ForksCount = repository.ForksCount
	existingRepo.StarsCount = repository.StarsCount
	existingRepo.OpenIssues = repository.OpenIssues
	existingRepo.Watchers = repository.Watchers
	existingRepo.CreatedAt = repository.CreatedAt
	existingRepo.UpdatedAt = repository.UpdatedAt
	existingRepo.SubscribersCount = repository.SubscribersCount
	return db.Save(&existingRepo).Error
}

func GetTopNCommitAuthors(db *gorm.DB, n, page, limit int) ([]struct {
	Author string `json:"author"`
	Count  int    `json:"count"`
}, error) {
	var authors []struct {
		Author string `json:"author"`
		Count  int    `json:"count"`
	}

	err := db.Model(&models.SavedCommit{}).
		Select("author, count(author) as count").
		Group("author").
		Order("count desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&authors).Error

	return authors, err
}
