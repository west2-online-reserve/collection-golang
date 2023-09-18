package dao

import (
	"context"
	"errors"
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

func (dao *TaskDao) SearchByText(uid uint, text string, start int) (tasks []*model.Task, count int64, err error) {
	text = "%" + text + "%"
	err = dao.DB.Model(&model.Task{}).Where("uid=? AND (content like ? OR title like ?)", uid, text, text).Count(&count).
		Limit(10).Offset((start - 1) * 10).Find(&tasks).Error
	return tasks, count, err
}

func (dao *TaskDao) SearchByStatus(uid uint, status int, start int) (tasks []*model.Task, count int64, err error) {
	err = dao.DB.Model(&model.Task{}).Where("uid=? AND status=?", uid, status).Count(&count).
		Limit(10).Offset((start - 1) * 10).Find(&tasks).Error
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

func (dao *TaskDao) DeleteAllTask(status int, uid uint) (count int64, err error) {
	switch status {
	case 0:
		err = dao.DB.Model(&model.Task{}).Where("status=?", status).Count(&count).Delete(&model.Task{}).Error
	case 1:
		err = dao.DB.Model(&model.Task{}).Where("status=?", status).Count(&count).Delete(&model.Task{}).Error
	case 3:
		err = dao.DB.Model(&model.Task{}).Where("uid=?", uid).Count(&count).Delete(&model.Task{}).Error
	default:
		return 0, errors.New("delete all task failed")
	}
	return count, err
}
