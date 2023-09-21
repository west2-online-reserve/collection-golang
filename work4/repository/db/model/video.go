package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Uid   uint
	Style int `gorm:"type:tinyint unsigned"`
	Url   string
}
