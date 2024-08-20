package repository

import (
	"fmt"
	"github-service/internal/models"

	"gorm.io/gorm"
)

func SaveCommit(db *gorm.DB, commit *models.SavedCommit) error {
	return db.Create(commit).Error
}

func GetCommitsByRepository(db *gorm.DB, repositoryURL string, page, limit int) ([]models.SavedCommit, error) {
	var commits []models.SavedCommit
	err := db.Where("url LIKE ?", fmt.Sprintf("%%%s%%", repositoryURL)).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&commits).Error
	return commits, err
}
