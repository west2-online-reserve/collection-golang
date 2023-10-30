package model

import (
	"gorm.io/gorm"
)

type Fans struct {
	gorm.Model
	Uid        uint
	FollowerId uint
}
