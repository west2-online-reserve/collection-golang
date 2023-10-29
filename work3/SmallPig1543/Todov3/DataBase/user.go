package DataBase

import (
	"Todov3/model"
)

func FindUserByUsername(userName string) (*model.User, error) {
	var user model.User
	err := DB.Model(&model.User{}).Where("username=?", userName).First(&user).Error
	return &user, err
}

func FindUserByUserID(id uint) (*model.User, error) {
	var user model.User
	err := DB.Model(&model.User{}).Where("id=?", id).First(&user).Error
	return &user, err
}

func CreateUser(user *model.User) error {
	err := DB.Model(&model.User{}).Create(user).Error
	return err
}
