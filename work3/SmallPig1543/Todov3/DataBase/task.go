package DataBase

import (
	"Todov3/model"
)

func CreateTask(task *model.Task) error {
	err := DB.Model(&model.Task{}).Create(task).Error
	return err
}

func FindTaskByUserIDAndID(uid, id uint) (*model.Task, error) {
	var task model.Task
	err := DB.Model(&model.Task{}).Where("id=? AND uid=?", id, uid).First(&task).Error
	return &task, err
}

func FindAllTasks(uid uint, limit, start int) ([]*model.Task, int64, error) {
	var list []*model.Task
	var count int64
	err := DB.Model(&model.Task{}).
		Preload("User").
		Where("uid=?", uid).
		Limit(limit).Offset((start - 1) * limit).
		Count(&count).Find(&list).Error
	return list, count, err
}

func FindAllTasksWithCondition(uid uint, condition, limit, start int) ([]*model.Task, int64, error) {
	var list []*model.Task
	var count int64
	err := DB.Model(&model.Task{}).
		Preload("User").
		Where("uid=?", uid).
		Where("status=?", condition).
		Limit(limit).Offset((start - 1) * limit).
		Count(&count).Find(&list).Error
	return list, count, err
}

func SearchTasks(uid uint, info string, limit, start int) ([]*model.Task, int64, error) {
	var list []*model.Task
	var count int64
	err := DB.Model(&model.Task{}).
		Where("uid=?", uid).
		Where("title LIKE ? OR content LIKE ?", "%"+info+"%", "%"+info+"%").
		Limit(limit).Offset((start - 1) * limit).
		Count(&count).
		Find(&list).Error
	return list, count, err
}

func UpdateTask(uid, id uint) error {
	task, err := FindTaskByUserIDAndID(uid, id)
	if err != nil {
		return err
	}
	if task.Status == 0 {
		task.Status = 1
	} else {
		task.Status = 0
	}
	return DB.Save(task).Error
}

func UpdateAllTasks(uid uint, condition int) error {
	var count int64
	DB.Model(&model.Task{}).Where("uid=?", uid).Count(&count)
	total := int(count)
	tasks, _, err := FindAllTasksWithCondition(uid, condition, total, 0)
	if err != nil {
		return err
	}
	for _, v := range tasks {
		if condition == 0 {
			v.Status = 1
		} else {
			v.Status = 0
		}
		err := DB.Save(v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteTask(uid, id uint) error {
	task, err := FindTaskByUserIDAndID(uid, id)
	if err != nil {
		return err
	}
	return DB.Delete(task).Error
}

func DeleteTasks(uid uint, condition int) error {
	var count int64
	DB.Model(&model.Task{}).Where("uid=?", uid).Count(&count)
	total := int(count)
	tasks, _, err := FindAllTasksWithCondition(uid, condition, total, 0)
	if err != nil {
		return err
	}
	for _, v := range tasks {
		if v.Status == condition {
			err := DB.Delete(v).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
