package DataBase

import (
	"Todov3/conf"
	"Todov3/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func LinkMySQL() {
	dsn := conf.MySqlUser + ":" + conf.MySqlPassword + "@tcp(127.0.0.1:3306)/" + conf.MySqlDataBase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}
	DB = db
	_ = DB.AutoMigrate(&model.User{})
	_ = DB.AutoMigrate(&model.Task{})
}
