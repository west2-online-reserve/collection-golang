package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"memogo/biz/dal/model"
	redisClient "memogo/biz/dal/redis"
	"time"

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

// 缓存键生成函数
func (r *UserRepository) userCacheKeyByID(id uint) string {
	return fmt.Sprintf("user:id:%d", id)
}

func (r *UserRepository) userCacheKeyByUsername(username string) string {
	return fmt.Sprintf("user:username:%s", username)
}

// 缓存失效：删除用户的所有缓存
func (r *UserRepository) invalidateUserCache(user *model.User) {
	if redisClient.RDB == nil {
		return
	}
	ctx := context.Background()

	// 删除 ID 缓存
	redisClient.RDB.Del(ctx, r.userCacheKeyByID(user.ID))

	// 删除用户名缓存
	if user.Username != "" {
		redisClient.RDB.Del(ctx, r.userCacheKeyByUsername(user.Username))
	}
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
	// 尝试从缓存获取
	if redisClient.RDB != nil {
		cacheKey := r.userCacheKeyByID(id)
		ctx := context.Background()

		cachedData, err := redisClient.RDB.Get(ctx, cacheKey).Result()
		if err == nil {
			// 缓存命中
			var user model.User
			if json.Unmarshal([]byte(cachedData), &user) == nil {
				return &user, nil
			}
		}
	}

	// 缓存未命中，查询数据库
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 写入缓存（10分钟过期，用户信息相对稳定）
	if redisClient.RDB != nil {
		cacheKey := r.userCacheKeyByID(id)
		ctx := context.Background()
		if data, err := json.Marshal(user); err == nil {
			redisClient.RDB.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	// 尝试从缓存获取
	if redisClient.RDB != nil {
		cacheKey := r.userCacheKeyByUsername(username)
		ctx := context.Background()

		cachedData, err := redisClient.RDB.Get(ctx, cacheKey).Result()
		if err == nil {
			// 缓存命中
			var user model.User
			if json.Unmarshal([]byte(cachedData), &user) == nil {
				return &user, nil
			}
		}
	}

	// 缓存未命中，查询数据库
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 写入缓存（10分钟过期）
	if redisClient.RDB != nil {
		cacheKey := r.userCacheKeyByUsername(username)
		ctx := context.Background()
		if data, err := json.Marshal(user); err == nil {
			redisClient.RDB.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	// 清除该用户的缓存
	r.invalidateUserCache(user)
	return nil
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id uint) error {
	// 先获取用户信息，用于清除缓存
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	// 清除该用户的缓存
	r.invalidateUserCache(user)
	return nil
}
