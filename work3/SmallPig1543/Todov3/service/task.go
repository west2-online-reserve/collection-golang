package service

import (
	"Todov3/DataBase"
	"Todov3/model"
	"Todov3/response"
	"Todov3/types"
	"context"
	"time"
)

type TaskService struct {
}

func (s *TaskService) Create(c context.Context, req *types.CreateTaskRequest) (interface{}, error) {
	//通过JWT从context取出UserID
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	//再通过UserID从数据库中获取用户信息
	user, err := DataBase.FindUserByUserID(u.ID)
	if err != nil {
		return response.BadResponse("数据库中无该用户"), err
	}
	task := model.Task{
		User:      *user,
		Uid:       user.ID,
		Title:     req.Title,
		Status:    req.Status,
		Content:   req.Content,
		StartTime: time.Now().Unix(),
	}
	err = DataBase.CreateTask(&task)
	if err != nil {
		return response.BadResponse("存入数据库失败"), err
	}
	return response.SuccessResponse(), nil
}

func (s *TaskService) ShowAllTasks(c context.Context, req *types.ShowAllTaskRequest) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	tasks, count, err := DataBase.FindAllTasks(u.ID, req.Limit, req.Start)
	if err != nil {
		return response.BadResponse("数据库出错"), err
	}
	taskList := make([]*response.TaskResp, 0)
	for _, v := range tasks {
		taskList = append(taskList, &response.TaskResp{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			View:      v.View(),
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Unix(),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
		v.AddView()
	}
	return response.TaskListRep(taskList, count), nil
}

func (s *TaskService) ShowAllTasksWithCondition(c context.Context, req *types.ShowAllTaskRequestWithCondition) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	tasks, count, err := DataBase.FindAllTasksWithCondition(u.ID, req.Condition, req.Limit, req.Start)
	if err != nil {
		return response.BadResponse("数据库出错"), err
	}
	taskList := make([]*response.TaskResp, 0)
	for _, v := range tasks {
		taskList = append(taskList, &response.TaskResp{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			View:      v.View(),
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Unix(),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
		v.AddView()
	}
	return response.TaskListRep(taskList, count), nil
}

func (s *TaskService) SearchTasks(c context.Context, req *types.SearchTasksRequest) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	tasks, count, err := DataBase.SearchTasks(u.ID, req.Info, req.Limit, req.Start)
	if err != nil {
		return response.BadResponse("数据库出错"), err
	}
	taskList := make([]*response.TaskResp, 0)
	for _, v := range tasks {
		taskList = append(taskList, &response.TaskResp{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			View:      v.View(),
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Unix(),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
		v.AddView()
	}
	return response.TaskListRep(taskList, count), nil
}

func (s *TaskService) UpdateTask(c context.Context, req *types.UpdateTaskRequest) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	err = DataBase.UpdateTask(u.ID, req.ID)
	if err != nil {
		return response.BadResponse("修改失败"), err
	}
	return response.SuccessResponse(), nil
}

func (s *TaskService) UpdateAllTasks(c context.Context, req *types.UpdateALlTasksRequest) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	err = DataBase.UpdateAllTasks(u.ID, req.Condition)
	if err != nil {
		return response.BadResponse("修改失败"), err
	}
	return response.SuccessResponse(), nil
}

func (s *TaskService) DeleteTask(c context.Context, req *types.DeleteTaskRequest) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	err = DataBase.DeleteTask(u.ID, req.ID)
	if err != nil {
		return response.BadResponse("删除失败"), err
	}
	return response.SuccessResponse(), nil
}

func (s *TaskService) DeleteTasks(c context.Context, req *types.DeleteTasksRequestWithCondition) (interface{}, error) {
	u, err := types.GetUserInfo(c)
	if err != nil {
		return response.BadResponse("获取用户信息失败"), err
	}
	err = DataBase.DeleteTasks(u.ID, req.Condition)
	if err != nil {
		return response.BadResponse("数据库出错"), err
	}
	return response.SuccessResponse(), nil
}
