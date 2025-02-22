package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string          `gorm:"size:255;not null;unique" json:"username" validate:"required,min=3,max=255"`
	Password     string          `gorm:"size:255;" json:"password"`
	RefreshToken string          `gorm:"size:255;" json:"refresh_token" `
	Provider     string          `gorm:"size:255;" json:"provider" validate:"omitempty,oneof=google facebook"`
	ProviderID   string          `gorm:"size:255;" json:"provider_id" validate:"omitempty,min=1"`
	Email        string          `gorm:"size:255;not null;unique" json:"email" validate:"required,email"`
	AvatarUrl    string          `gorm:"size:255" json:"avatar_url" validate:"omitempty,url"`
	Name         string          `gorm:"size:255;not null" json:"name" validate:"required,min=1,max=255"`
	Role         string          `gorm:"size:255;default:user" json:"role" validate:"omitempty,oneof=user admin"`
	Verified     bool            `gorm:"default:false" json:"verified"`
	LastLogin    *gorm.DeletedAt `gorm:"index" json:"last_login"`
}

// Fungsi untuk melakukan validasi pada model User
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
