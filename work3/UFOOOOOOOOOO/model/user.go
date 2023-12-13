package model

import (
	// "github.com/henrylee2cn/goutil/password"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	PasswordDigest string
}


//加密密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

//验证密码
func  (user *User) CheckPassword (password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), 
	[]byte(password))
	return err==nil
}