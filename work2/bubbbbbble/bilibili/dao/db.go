package dao

import (
	"bilibili/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var DB *gorm.DB
func InitMysql(){
	dsn := "root:123456@tcp(127.0.0.1:3306)/mydata?charset=utf8mb4&parseTime=True&loc=Local"
	db,_:=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	DB=db
	DB.AutoMigrate(&model.StoredReply{})
}