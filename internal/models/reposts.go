package models

import (
	"gorm.io/gorm"
)

type Repost struct {
	gorm.Model

	UserID    uint   `gorm:"not null" json:"user_id"`
	BlogID    uint   `gorm:"not null" json:"blog_id"`
	Quote     string `gorm:"type:text" json:"quote"`
	LikeCount int    `gorm:"type:integer;not null;default:0" json:"like_count"`

	// Relation
	Owner *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"owner"`
	Blog  *Blog `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE" json:"blog"`
}
