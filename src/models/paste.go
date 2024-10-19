package models

import (
	"gorm.io/gorm"
)

type Paste struct {
	gorm.Model
	Title     string
	Content   string
	Extension string
	UserID    uint
	Slug      string `gorm:"unique"`
}
