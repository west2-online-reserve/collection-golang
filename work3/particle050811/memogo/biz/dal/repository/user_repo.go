package repository

import (
	"errors"
	"memogo/biz/dal/model"

	"gorm.io/gorm"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	// 检查用户名是否已存在
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrUserAlreadyExists
	}

	return r.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
