package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Uid     uint
	Content string `gorm:"type:text"`
}
