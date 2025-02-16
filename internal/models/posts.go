package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model // Ini akan menambahkan ID, CreatedAt, UpdatedAt, DeletedAt

	Title    string         `gorm:"size:255;not null" json:"title"`
	Body     string         `gorm:"type:text;not null" json:"body"`
	ReferID  *uint          `gorm:"index" json:"refer_id"`
	OwnerID  uint           `gorm:"not null" json:"owner_id"`
	Category pq.StringArray `gorm:"type:text[]" json:"category"` // Array PostgreSQL

	// Relation
	Refer *Post `gorm:"foreignKey:ReferID;constraint:OnDelete:SET NULL" json:"refer"`
	Owner *User `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner"`
}
