package service

import (
	"hertz/pkg/model"
	"hertz/pkg/repository"
)

type todoService struct {
	tr repository.TodoRepository
}

type TodoService interface {
	AddTodo(todo *model.Todo) bool
	SetTodoFinishedById(uid uint64, id uint64) bool
	SetTodoNotFinishedById(uid uint64, id uint64) bool
	SetTodosFinished(uid uint64) bool
	SetTodosNotFinished(uid uint64) bool
	GetFinishedTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64)
	GetNotFinishedTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64)
	GetTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64)
	GetTodosByKey(uid uint64, key string, pageNum, pageSize int) ([]*model.Todo, int64)
	DeleteTodoById(id uint64, uid uint64) bool
	DeleteFinishedTodo(uid uint64) bool
	DeleteNotFinishedTodo(uid uint64) bool
	DeleteAll(uid uint64) bool
}

func NewTodoService(tr repository.TodoRepository) TodoService {
	return &todoService{tr: tr}
}

func (ts *todoService) AddTodo(todo *model.Todo) bool {
	return ts.tr.Create(todo)
}

func (ts *todoService) SetTodoFinishedById(uid uint64, id uint64) bool {
	return ts.tr.UpdateStatusById(uid, id, 1)
}

func (ts *todoService) SetTodoNotFinishedById(uid uint64, id uint64) bool {
	return ts.tr.UpdateStatusById(uid, id, 0)
}

func (ts *todoService) SetTodosFinished(uid uint64) bool {
	return ts.tr.UpdateFinishedStatus(uid)
}

func (ts *todoService) SetTodosNotFinished(uid uint64) bool {
	return ts.tr.UpdateNotFinishedStatus(uid)
}

func (ts *todoService) GetFinishedTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64) {
	return ts.tr.QueryByStatus(uid, 1, pageNum, pageSize)
}

func (ts *todoService) GetNotFinishedTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64) {
	return ts.tr.QueryByStatus(uid, 0, pageNum, pageSize)
}

func (ts *todoService) GetTodos(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64) {
	return ts.tr.Query(uid, pageNum, pageSize)
}

func (ts *todoService) GetTodosByKey(uid uint64, key string, pageNum, pageSize int) ([]*model.Todo, int64) {
	return ts.tr.QueryByKey(uid, key, pageNum, pageSize)
}

func (ts *todoService) DeleteTodoById(id uint64, uid uint64) bool {
	return ts.tr.DeleteById(id, uid)
}

func (ts *todoService) DeleteAll(uid uint64) bool {
	return ts.tr.DeleteAll(uid)
}

func (ts *todoService) DeleteFinishedTodo(uid uint64) bool {
	return ts.tr.DeleteByStatus(uid, 1)
}

func (ts *todoService) DeleteNotFinishedTodo(uid uint64) bool {
	return ts.tr.DeleteByStatus(uid, 0)
}
