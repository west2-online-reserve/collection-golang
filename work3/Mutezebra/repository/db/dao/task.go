package dao

import (
	"context"
	"gorm.io/gorm"
	"three/repository/db/model"
)

type TaskDao struct {
	*gorm.DB
}

// NewTaskDao return a pointer to TaskDao with DB info
func NewTaskDao(ctx context.Context) *TaskDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &TaskDao{NewDBClient(ctx)}
}

func (dao *TaskDao) FindTaskById(id uint, uid uint) (task *model.Task, err error) {
	err = dao.DB.Model(&model.Task{}).Where("id=? AND uid=?", id, uid).Find(&task).Error
	return
}

func (dao *TaskDao) Create(task *model.Task) error {
	err := dao.DB.Model(&model.Task{}).Create(task).Error
	return err
}

func (dao *TaskDao) Update(task *model.Task) error {
	err := dao.DB.Save(task).Error
	return err
}

func (dao *TaskDao) SearchByText(uid uint, text string) (tasks []*model.Task, count int64, err error) {
	text = "%" + text + "%"
	err = dao.DB.Model(&model.Task{}).Where("uid=? AND (content like ? OR title like ?)", uid, text, text).Find(&tasks).Count(&count).Error
	return tasks, count, err
}

func (dao *TaskDao) SearchAll(uid uint, status int) (tasks []*model.Task, count int64, err error) {
	err = dao.DB.Model(&model.Task{}).Where("uid=? AND status=?", uid, status).Find(&tasks).Count(&count).Error
	return tasks, count, err
}

func (dao *TaskDao) List(uid uint, limit int, start int) (tasks []*model.Task, count int64, err error) {
	err = dao.DB.Model(&model.Task{}).Where("uid=?", uid).Count(&count).
		Limit(limit).Offset((start - 1) * limit).
		Find(&tasks).Error
	return tasks, count, err
}

func (dao *TaskDao) Delete(task *model.Task) error {
	err := dao.DB.Delete(&model.Task{}, task).Error
	return err
}
