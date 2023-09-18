package dao

import (
	"context"
	"gorm.io/gorm"
	"three/repository/db/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUserName(userName string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).
		Where("user_name=?", userName).
		First(&user).
		Error
	return user, err
}

func (dao *UserDao) Create(user *model.User) error {
	err := dao.DB.Model(&model.User{}).
		Create(&user).
		Error
	return err
}
