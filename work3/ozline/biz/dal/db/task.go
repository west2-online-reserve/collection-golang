package db

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants"
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
		Table(constants.TaskTable).
		Create(&taskResp).
		Error

	if err != nil {
		return nil, err
	}

	return taskResp, nil
}

func UpdateTask(ctx context.Context, userid, id, status int64) (err error) {

	if id == 0 {
		err = DB.WithContext(ctx).
			Table(constants.TaskTable).
			Where("user_id = ?", userid).
			Update("status", status).
			Error
	} else {
		err = DB.WithContext(ctx).
			Table(constants.TaskTable).
			Where("user_id = ?", userid).
			Where("id = ?", id).
			Update("status", status).
			Error
	}

	return
}

func DeleteTaskSingle(ctx context.Context, userid, id int64) error {
	hlog.CtxInfof(ctx, "useid: %v\n", userid)
	return DB.WithContext(ctx).
		Table(constants.TaskTable).
		Where("user_id = ? AND id = ?", userid, id).
		Delete(&Task{}).
		Error
}

func DeleteTask(ctx context.Context, userid, status int64) (err error) {

	if status == 2 { // 删除全部内容
		err = DB.WithContext(ctx).
			Table(constants.TaskTable).
			Delete(&Task{UserId: userid}).
			Error
	}

	if status == 1 { // 删除已完成内容
		err = DB.WithContext(ctx).
			Table(constants.TaskTable).
			Where("status = ?", 1).
			Delete(&Task{UserId: userid}).
			Error
	}

	if status == 0 { // 删除已完成内容
		err = DB.WithContext(ctx).
			Table(constants.TaskTable).
			Where("status = ?", 0).
			Delete(&Task{UserId: userid}).
			Error
	}

	return err
}

func QueryTaskListByStatus(ctx context.Context, userid, pasgesize, pagenum, status int64) ([]*Task, int64, error) {
	var taskResp []*Task
	var count int64

	err := DB.WithContext(ctx).
		Table(constants.TaskTable).
		Where("status = ?", status).
		Where("user_id = ?", userid).
		Limit(int(pasgesize)).
		Offset(int((pagenum - 1) * pasgesize)).
		Count(&count).
		Find(&taskResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return taskResp, count, nil
}

func QueryTaskListByKeyword(ctx context.Context, userid, pasgesize, pagenum int64, keywords string) ([]*Task, int64, error) {
	var taskResp []*Task
	var count int64

	err := DB.WithContext(ctx).
		Table(constants.TaskTable).
		Where("title LIKE ?", keywords).
		Where("user_id = ?", userid).
		Limit(int(pasgesize)).
		Offset(int((pagenum - 1) * pasgesize)).
		Count(&count).
		Find(&taskResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return taskResp, count, nil
}
