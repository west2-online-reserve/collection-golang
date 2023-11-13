package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(connect string) {
	db, err := gorm.Open(mysql.Open(connect), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //开启日志
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数形式的表名
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("连接数据库失败")
		return
	}
	fmt.Println("连接数据库成功")
	err = db.Migrator().DropTable(&Article{})
	if err != nil {
		return
	}
	mysqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("获取数据库连接对象失败")
		return
	}
	mysqlDB.SetMaxIdleConns(1000)                // 设置最大空闲连接数
	mysqlDB.SetMaxOpenConns(1000)                // 设置最大连接数
	mysqlDB.SetConnMaxLifetime(20 * time.Second) // 设置连接的最大可复用时间
	DB = db
	err = db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&Article{})
	if err != nil {
		fmt.Println("数据库建表失败")
		return
	}
}

// InsertArticle 插入数据
func InsertArticle(article *Article, db *gorm.DB) error {
	tx := db.Begin()
	result := tx.Create(article)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("数据插入失败")
		return result.Error
	}
	tx.Commit()
	return nil
}
