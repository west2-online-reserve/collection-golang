package db

import (
	"Demo/pkg/constants"
	"context"
	"time"
)

func CreateTask(ctx context.Context, title, content string, userid, startAt, endAt int64) (*Task, error) {
	var taskResp *Task

	taskResp = &Task{
		UserId:  userid,
		Title:   title,
		Content: content,
		StartAt: time.Unix(startAt, 0),
		EndAt:   time.Unix(endAt, 0),
	}

	err := DB.WithContext(ctx).
		Table(constants.TableTask).
		Create(&taskResp).
		Error

	if err != nil {
		return nil, err
	}
	return taskResp, nil
}

func UpdateTask(ctx context.Context, userid, id, status int64) error {
	var err error
	if id == 0 {
		err = DB.WithContext(ctx).
			Table(constants.TableTask).
			Where("user_id = ?", userid).
			Update("status", status).
			Error
	} else {
		err = DB.WithContext(ctx).
			Table(constants.TableTask).
			Where("user_id = ?", userid).
			Where("id=?", id).
			Update("status", status).
			Error
	}
	return err
}

func DeleteTaskSingle(ctx context.Context, userid, id int64) error {

	return DB.WithContext(ctx).
		Table(constants.TableTask).
		Where("id = ?", id).
		Delete(&Task{Id: id, UserId: userid}).
		Error

}

func DeleteTask(ctx context.Context, userid, status int64) (err error) {
	if status == 2 { //删除全部内容
		err = DB.WithContext(ctx).
			Table(constants.TableTask).
			Delete(&Task{UserId: userid}).
			Error
	}

	if status == 1 { //删除已经完成的内容
		err = DB.WithContext(ctx).
			Table(constants.TableTask).
			Where("status = ?", 1).
			Delete(&Task{UserId: userid}).
			Error
	}

	if status == 0 { //删除还未完成的内容
		err = DB.WithContext(ctx).
			Table(constants.TableTask).
			Where("status = ?", 0).
			Delete(&Task{UserId: userid}).
			Error
	}
	return err
}

func QuerySingleTask(ctx context.Context, userid, id int64) (*Task, error) {
	var taskResp *Task
	err := DB.WithContext(ctx).
		Table(constants.TableTask).
		Where("user_id = ?", userid).
		Where("id=?", id).
		Where("deleted_at IS NULL"). // 排除已软删除的记录
		First(&taskResp).            //应该用first否则查找不到会返回空的response
		Error
	if err != nil {
		return nil, err
	}
	return taskResp, nil
}

func QueryTaskListByStatus(ctx context.Context, userid, pagesize, pagenum,
	status int64) ([]*Task, int64, error) {
	var taskResp []*Task
	var count int64

	err := DB.WithContext(ctx).
		Table(constants.TableTask).
		Where("status = ?", status).
		Where("user_id = ?", userid).
		Where("deleted_at IS NULL"). // 排除已软删除的记录
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&taskResp).
		Error
	if err != nil {
		return nil, -1, err
	}
	return taskResp, count, nil
}

func QueryTaskListByKeyword(ctx context.Context, userid, pagesize, pagenum int64,
	keyword string) ([]*Task, int64, error) {
	var taskResp []*Task
	var count int64

	err := DB.WithContext(ctx).
		Table(constants.TableTask).
		Where("title LIKE ?", keyword).
		Where("user_id = ?", userid).
		Where("deleted_at IS NULL"). // 排除已软删除的记录
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&taskResp).
		Error
	if err != nil {
		return nil, -1, err
	}
	return taskResp, count, nil
}
