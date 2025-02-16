package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Embed gorm.Model untuk field umum
	GoogleId   string `gorm:"size:255;not null" json:"google_id"`
	Email      string `gorm:"size:255;not null;unique" json:"email"` // Tambahkan unique
	AvatarUrl  string `gorm:"size:255;not null" json:"avatar_url"`
	Name       string `gorm:"size:255;not null" json:"name"`
}
