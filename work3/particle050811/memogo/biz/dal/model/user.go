package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username     string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	PasswordHash string `gorm:"not null;size:255" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
