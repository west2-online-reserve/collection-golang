package repository

import (
	"hertz/pkg/model"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	QueryUser(u *model.User) bool
	Register(u *model.User) bool
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) QueryUser(u *model.User) bool {
	if err := ur.db.First(u).Error; err != nil {
		log.Fatalf("用户查询失败: %v", err)
		return false
	}
	return true
}

func (ur *userRepository) Register(u *model.User) bool {
	if err := ur.db.Create(u).Error; err != nil {
		log.Fatalf("创建用户失败: %v", err)
		return false
	}
	return true
}
