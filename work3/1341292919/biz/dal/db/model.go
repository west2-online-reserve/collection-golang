package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Task struct {
	Id        int64
	UserId    int64
	Title     string
	Content   string
	Status    int64
	StartAt   time.Time
	EndAt     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
