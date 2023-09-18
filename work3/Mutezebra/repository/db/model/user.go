package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"three/consts"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
	NickName string
}

func (*User) TableName() string {
	return "user"
}

func (u *User) SetPassword(password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), consts.PasswordCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}
