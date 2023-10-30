package dao

import (
	"four/pkg/log"
	"four/repository/db/model"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Video{})
	if err != nil {
		log.LogrusObj.Error(err)
	}
}
