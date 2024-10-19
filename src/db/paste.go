package db

import (
	"gopastebin/src/models"
)

func CreatePaste(paste *models.Paste) error {
	return db.Create(paste).Error
}

func GetPasteBySlug(slug string, paste *models.Paste) error {
	return db.Where("slug = ?", slug).First(paste).Error
}

func UpdatePaste(paste *models.Paste) error {
	return db.Save(paste).Error
}

func IsPasteUnique(slug string) (bool, error) {
	var count int64
	err := db.Model(&models.Paste{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
