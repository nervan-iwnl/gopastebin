package service

import (
	"strings"

	"gopastebin/internal/repository"
)

const settingStorageProvider = "storage_provider"

type AppSettingsService struct {
	repo           *repository.SettingRepository
	defaultStorage string
}

func NewAppSettingsService(repo *repository.SettingRepository, defaultStorage string) *AppSettingsService {
	return &AppSettingsService{
		repo:           repo,
		defaultStorage: defaultStorage,
	}
}

func (s *AppSettingsService) GetStorageProvider() string {
	val, err := s.repo.Get(settingStorageProvider)
	if err != nil || val == "" {
		return s.defaultStorage
	}
	return val
}

func (s *AppSettingsService) SetStorageProvider(v string) error {
	v = strings.ToLower(v)
	if v != "firebase" && v != "local" {
		return nil
	}
	return s.repo.Set(settingStorageProvider, v)
}
