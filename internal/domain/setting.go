package domain

import "time"

// Setting — простая k/v настройка, чтобы хранить текущий storage_provider
type Setting struct {
	Key       string    `gorm:"primaryKey;size:100"`
	Value     string    `gorm:"size:500"`
	UpdatedAt time.Time
}
