package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
	Author string
	Title  string
	Num    int
	Time   string
	Text   string
}

var db *gorm.DB

func Init() {
	//"root:9055@tcp(127.0.0.1:3306)/test_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:9055@tcp(127.0.0.1:3306)/info?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//创建表 自动迁移(将结构体与数据表对应起来)
	db = d
}
func CreatTable() {
	//创建表 自动迁移(将结构体与数据表对应起来)
	err := db.AutoMigrate(&Data{})
	if err != nil {
		panic(err)
	}
}
func insert(data *Data) {
	db.Create(&data)
}
