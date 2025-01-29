package models

type Post struct {
	Id    int64  `gorm:"primaryKey" json:"id"`
	Title string `gorm:"type:varchar(100)" json:"title"`
	Body  string `gorm:"type:varchar(300)" json:"body"`
}
