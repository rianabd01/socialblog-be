package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model

	UserID *uint `gorm:"not null;uniqueIndex:idx_user_blog_post" json:"user_id"`
	BlogID *uint `gorm:"uniqueIndex:idx_user_blog_post" json:"blog_id,omitempty"` // Bisa null
	PostID *uint `gorm:"uniqueIndex:idx_user_blog_post" json:"post_id,omitempty"` // Bisa null

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user"`
	Blog *Blog `gorm:"foreignKey:BlogID" json:"blog"`
	Post *Post `gorm:"foreignKey:PostID" json:"post"`
}
