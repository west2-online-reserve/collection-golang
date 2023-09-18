package service

import (
	"context"
	"fmt"
	"sync"
	"three/pkg/ctl"
	"three/pkg/e"
	"three/repository/db/dao"
	"three/repository/db/model"
	"three/types"
	"time"
)

type TaskService struct {
}

var taskSrvOnce sync.Once
var TaskSrvIns *TaskService

func GetTaskSrv() *TaskService {
	if TaskSrvIns != nil {
		return TaskSrvIns
	}
	taskSrvOnce.Do(func() {
		TaskSrvIns = &TaskService{}
	})
	return TaskSrvIns
}

func (s *TaskService) Create(ctx context.Context, req *types.TaskCreateReq) (resp interface{}, err error) {
	var code int
	code = e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	task := &model.Task{
		Uid:     userInfo.ID,
		Title:   req.Title,
		Content: req.Content,
		Status:  0,
	}

	task.StartTime = time.Now().Unix()
	err = taskDao.Create(task)
	resp = &types.TaskInfoResp{
		Id:      task.ID,
		Title:   req.Title,
		Content: req.Content,
		Status:  0,
	}
	if err != nil {
		code = e.CreateTaskFailed
		return ctl.RespError(err, code), err
	}
	return ctl.RespSuccessWithData(resp, code), nil
}

func (s *TaskService) Update(ctx context.Context, req *types.TaskUpdateReq) (interface{}, error) {
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	task, err := taskDao.FindTaskById(req.Id, userInfo.ID)
	if err != nil {
		code = e.FindTaskFailed
		return ctl.RespError(err, code), err
	}
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Content != "" {
		task.Content = req.Content
	}
	if req.Status == 1 {
		task.EndTime = time.Now().Unix()
	}
	task.Status = req.Status

	err = taskDao.Update(task)
	if err != nil {
		code = e.UpdateTaskFailed
		return ctl.RespError(err, code), err
	}
	taskInfo := &types.TaskInfoResp{
		Title:     task.Title,
		Content:   task.Content,
		Status:    task.Status,
		View:      task.View(),
		CreateAt:  task.CreatedAt.Unix(),
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
	}
	return ctl.RespSuccessWithData(taskInfo, code), nil
}

func (s *TaskService) Show(ctx context.Context, req *types.TaskShowReq) (interface{}, error) {
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	task, err := taskDao.FindTaskById(req.Id, userInfo.ID)
	if err != nil {
		code = e.FindTaskFailed
		return ctl.RespError(err, code), err
	}

	taskInfo := &types.TaskInfoResp{
		Title:     task.Title,
		Content:   task.Content,
		Status:    task.Status,
		View:      task.View(),
		CreateAt:  task.CreatedAt.Unix(),
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
	}
	task.AddView()
	return ctl.RespSuccessWithData(taskInfo, code), nil
}

func (s *TaskService) List(ctx context.Context, req *types.TaskListReq) (interface{}, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	tasks, count, err := taskDao.List(userInfo.ID, req.Limit, req.Start)
	if err != nil {
		code = e.ListTaskFailed
		return ctl.RespError(err, code), err
	}
	results := make([]*types.TaskInfoResp, 0)
	for _, item := range tasks {
		data := &types.TaskInfoResp{
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			View:      item.View(),
			CreateAt:  item.CreatedAt.Unix(),
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
		}
		results = append(results, data)
	}
	return ctl.RespSuccessWithData(ctl.TaskItemList{Count: count, Items: results, Page: req.Start + 1}, code), nil
}

func (s *TaskService) SearchByText(ctx context.Context, req *types.TaskSearchReq) (resp interface{}, err error) {
	if req.Start == 0 {
		req.Start = 1
	}
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	tasks, count, err := taskDao.SearchByText(userInfo.ID, req.Text, req.Start)
	if err != nil {
		code = e.SearchTaskFailed
		return ctl.RespError(err, code), err
	}
	results := make([]*types.TaskInfoResp, 0)
	for _, item := range tasks {
		data := &types.TaskInfoResp{
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			View:      item.View(),
			CreateAt:  item.CreatedAt.Unix(),
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
		}
		results = append(results, data)
	}
	return ctl.RespSuccessWithData(ctl.TaskItemList{Count: count, Items: results, Page: req.Start}, code), nil
}

func (s *TaskService) SearchByStatus(ctx context.Context, req *types.TaskSearchReq) (interface{}, error) {
	if req.Start == 0 {
		req.Start = 1
	}
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	tasks, count, err := taskDao.SearchByStatus(userInfo.ID, req.Status, req.Start)
	if err != nil {
		code = e.SearchTaskFailed
		return ctl.RespError(err, code), err
	}
	results := make([]*types.TaskInfoResp, 0)
	for _, item := range tasks {
		data := &types.TaskInfoResp{
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			View:      item.View(),
			CreateAt:  item.CreatedAt.Unix(),
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
		}
		results = append(results, data)
	}
	return ctl.RespSuccessWithData(ctl.TaskItemList{Count: count, Items: results, Page: req.Start}, code), nil
}

func (s *TaskService) Delete(ctx context.Context, req *types.TaskDeleteReq) (interface{}, error) {
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}
	task, err := taskDao.FindTaskById(req.Id, userInfo.ID)
	if err != nil {
		code = e.FindTaskFailed
		return ctl.RespError(err, code), err
	}
	err = taskDao.Delete(task)
	if err != nil {
		code = e.DeleteTaskFailed
		return ctl.RespError(err, code), err
	}
	return ctl.RespSuccess(code), nil
}

func (s *TaskService) DeleteAllTask(ctx context.Context, req *types.TaskDeleteReq) (interface{}, error) {
	code := e.SUCCESS
	taskDao := dao.NewTaskDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil && (userInfo.ID == 0 || userInfo.UserName == "") {
		code = e.GetUserInfoFailed
		return ctl.RespError(err, code), err
	}

	count, err := taskDao.DeleteAllTask(req.Status, userInfo.ID)
	if err != nil {
		code = e.DeleteTaskFailed
		return ctl.RespError(err, code), err
	}
	return ctl.RespSuccessWithData(fmt.Sprintf("delete %d of tasks", count), code), nil
}
