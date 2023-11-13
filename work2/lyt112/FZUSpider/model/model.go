package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Time     string
	Title    string
	Content  string
	Author   string
	ClickNum int
}
