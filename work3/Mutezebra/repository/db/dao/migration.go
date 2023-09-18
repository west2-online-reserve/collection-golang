package dao

import (
	"three/pkg/utils"
	"three/repository/db/model"
)

// migration 数据迁移
func migration() {
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		utils.LogrusObj.Errorln(err)
	}
}
