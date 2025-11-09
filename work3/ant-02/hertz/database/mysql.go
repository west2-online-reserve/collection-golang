package database

import (
	"hertz/config"
	"hertz/pkg/model"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlOnce sync.Once
	db        *gorm.DB
)

func GetMysql() *gorm.DB {
	var err error = nil
	mysqlOnce.Do(func() {
		cfg := config.GetConfig()
		dsn := cfg.Database.Username + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ":" + cfg.Database.Port + ")/" + cfg.Database.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	})
	if err != nil {
		log.Fatalf("fail to get mysql: %v", err)
		return nil
	}
	autoMigrate()
	return db
}

func autoMigrate() {
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
}
