package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online-reserve/collection-golang/work3/biz/dal/db"
	"github.com/west2-online-reserve/collection-golang/work3/biz/model/task"
)

type TaskService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewTaskService(ctx context.Context, c *app.RequestContext) *TaskService {
	return &TaskService{ctx: ctx, c: c}
}

func (s *TaskService) AddTask(req *task.AddTaskRequest) (*db.Task, error) {
	return db.CreateTask(s.ctx, req.Title, req.Content, GetUserIDFromContext(s.c), req.StartAt, req.EndAt)
}

func (s *TaskService) UpdateTask(req *task.UpdateTaskRequest) error {
	return db.UpdateTask(s.ctx, GetUserIDFromContext(s.c), req.ID, req.Status)
}

func (s *TaskService) DeleteTaskSingle(req *task.DeleteTaskRequest) error {
	return db.DeleteTaskSingle(s.ctx, GetUserIDFromContext(s.c), req.ID)
}

func (s *TaskService) DeleteTaskByStatus(req *task.DeleteTaskByStatusRequest) error {
	return db.DeleteTask(s.ctx, GetUserIDFromContext(s.c), req.Status)
}

func (s *TaskService) QueryTaskListByStatus(req *task.QueryTaskListByStatusRequest) ([]*db.Task, int64, error) {
	return db.QueryTaskListByStatus(s.ctx, GetUserIDFromContext(s.c), req.PageSize, req.PageNum, req.Status)
}

func (s *TaskService) QueryTaskListByKeyword(req *task.QueryTaskListByKeywordRequest) ([]*db.Task, int64, error) {
	return db.QueryTaskListByKeyword(s.ctx, GetUserIDFromContext(s.c), req.PageSize, req.PageNum, req.Keyword)
}
