package dao

import (
	"bilibili/db/model"
	"gorm.io/gorm"
)

type InfoDao struct {
	*gorm.DB
}

func (dao *InfoDao) Add(info *model.Info) error {
	if err := dao.DB.Model(&model.Info{}).Create(&info).Error; err != nil {
		return err
	}
	return nil
}
