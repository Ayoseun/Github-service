package repository

import (
	"github-service/internal/models"

	"gorm.io/gorm"
)

func SaveRepository(db *gorm.DB, repository *models.Repository) error {
	return db.Create(repository).Error
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
