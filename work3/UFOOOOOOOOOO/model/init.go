package model

import (
	"fmt"
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func Database(connstring string){
	db, err := gorm.Open("mysql", connstring)
	if err != nil {
		fmt.Println(err)
		panic("Mysql数据库连接错误")
	}
	fmt.Println("数据库连接成功")
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(20)// 设置连接池
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second*30)
	DB=db
	migration()
}