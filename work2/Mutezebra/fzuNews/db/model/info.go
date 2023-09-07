package model

type Info struct {
	Title   string `gorm:"index;not null"`
	Content string `gorm:"type:longtext"`
	Author  string
	Date    string
	Views   string
}
