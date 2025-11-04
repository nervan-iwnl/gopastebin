package repository

import (
	"gopastebin/internal/domain"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) Get(key string) (string, error) {
	var s domain.Setting
	if err := r.db.First(&s, "key = ?", key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return s.Value, nil
}

func (r *SettingRepository) Set(key, value string) error {
	s := domain.Setting{
		Key:   key,
		Value: value,
	}
	return r.db.Save(&s).Error
}
