package repository

import (
	"gopastebin/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Paste{})
}
