package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model

	UserID    uint           `gorm:"not null" json:"user_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	ImageUrl  string         `gorm:"size:255;not null" json:"image_url"`
	Category  pq.StringArray `gorm:"type:text[]" json:"category"`
	LikeCount int            `gorm:"type:integer;not null;default:0" json:"like_count"`

	// Relation
	Owner *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"owner"`
}
