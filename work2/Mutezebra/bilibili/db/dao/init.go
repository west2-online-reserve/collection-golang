package dao

import (
	"bilibili/conf"
	"bilibili/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"time"
)

var _db *gorm.DB

// MysqlInit 初始化数据库
func MysqlInit() {
	mConfig := conf.Config.Mysql["default"]

	dsn := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":",
		mConfig.DbPort, ")/", mConfig.DbName, "?charset=", mConfig.Charset, "&parseTime=True&loc=Local"}, "")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(15 * time.Second)
	_db = db
	migration()
	log.Println("数据库连接成功")
}

// 自动迁移
func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.Info{})
	if err != nil {
		return
	}
}

func Db() *gorm.DB {
	return _db
}
