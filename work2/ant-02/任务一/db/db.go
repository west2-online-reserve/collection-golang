package db

import (
	"log"
	"os"
	"sync"
	"time"
	"west2/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type niDb struct {
	db      *gorm.DB
	dbMutex sync.Mutex
}

var (
	once sync.Once
	nd   *niDb
)

func InitDb() *niDb {
	once.Do(func() {
		nd = &niDb{}
		dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		nd.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("连接数据库失败: %v", err)
			return
		}

		sqlDB, err := nd.db.DB()
		if err != nil {
			log.Fatalf("获取数据库实例失败: %v", err)
			return
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		if err := nd.db.AutoMigrate(&model.NotiInfo{}); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
			return
		}

	})

	return nd
}

func (nd *niDb) AddNi(ni *model.NotiInfo) {
	if ni == nil {
		return
	}
	nd.dbMutex.Lock()
	defer nd.dbMutex.Unlock()
	if err := nd.db.Create(ni).Error; err != nil {
		log.Fatalf("failed to add ni: %v\nerr: %v", ni, err)
		return
	}
}
