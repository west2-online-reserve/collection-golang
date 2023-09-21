package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	PasswordDigest string
	Follow         int
	Fans           int
	NickName       string
	Avatar         string
	Gender         int `gorm:"type:tinyint unsigned"`
	Email          string
}
