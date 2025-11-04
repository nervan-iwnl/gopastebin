package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email           string `gorm:"uniqueIndex" json:"email"`
	Username        string `gorm:"uniqueIndex" json:"username"`
	Password        string `json:"-"`
	Avatar          string `json:"avatar"`
	EmailVerified   bool   `json:"email_verified"`
	EmailVerifyCode string `json:"-"`
	IsAdmin       bool   `gorm:"default:false" json:"is_admin"` // ← добавили
}

func ExtractUserID(v any) uint {
	if v == nil {
		return 0
	}
	if u, ok := v.(*User); ok {
		return u.ID
	}
	return 0
}
