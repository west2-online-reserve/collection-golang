package model

import (
	"fmt"
	"fzu_crawl/conf"
	"fzu_crawl/data/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitDB() *gorm.DB {
	s := conf.InitConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", s.Username, s.Password, s.Host, s.Port, s.Database)
	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	migErr := db.AutoMigrate(&model.News{})
	if migErr != nil {
		panic("failed to migrate")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	return db
}

func CloseDB(db *gorm.DB) error {
	DB, err := db.DB()
	if err != nil {
		return err
	}
	return DB.Close()
}
