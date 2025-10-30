package db

import (
	"log"

	"memogo/biz/dal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接（使用 SQLite）
func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("memogo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database initialized successfully")
}
