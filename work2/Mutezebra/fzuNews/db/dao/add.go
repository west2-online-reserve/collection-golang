package dao

import (
	model2 "fzuNews/db/model"
	"gorm.io/gorm"
)

type InfoDao struct {
	*gorm.DB
}

func (dao *InfoDao) Add(info model2.Info) error {
	if err := dao.DB.Model(&model2.Info{}).Create(&info).Error; err != nil {
		return err
	}
	return nil
}
