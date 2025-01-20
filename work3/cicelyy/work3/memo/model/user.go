//建表
package model

import (
	"github.com/jinzhu/gorm" //用于简化数据库操作
	"golang.org/x/crypto/bcrypt" //用于安全地存储和验证密码
)

//定义结构体User
type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	PasswordDigest string //存储加密后的密码
}

//加密密码
func (user *User) SetPassword(password string) error {
	//使用 bcrypt.GenerateFromPassword 函数生成密码的哈希值
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err 
	}
	//将生成的哈希值转换为字符串并存储
	user.PasswordDigest = string(bytes)
	return nil
}

//验证密码
func  (user *User) CheckPassword (password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest),[]byte(password))  
	return err==nil
} 