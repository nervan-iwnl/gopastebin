package domain

import (
	"time"

	"gorm.io/gorm"
)

type Paste struct {
	gorm.Model
	Title       string     `json:"title"`
	Slug        string     `gorm:"uniqueIndex" json:"slug"`
	StoragePath string     `json:"storage_path"` // реальный путь в Firebase/local
	Extension   string     `json:"extension"`
	UserID      uint       `json:"user_id"`
	Folder      string     `json:"folder"` // "project-euler/001"
	IsPublic    bool       `json:"is_public"`
	ExpiresAt   *time.Time `json:"expires_at"`
}
