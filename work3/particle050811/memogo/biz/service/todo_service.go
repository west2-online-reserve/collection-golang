package service

import (
    "errors"
    "time"

    "memogo/biz/dal/model"
    "memogo/biz/dal/repository"
)

var (
    // ErrTitleRequired 标题必填
    ErrTitleRequired = errors.New("title is required")
    // ErrContentRequired 内容必填
    ErrContentRequired = errors.New("content is required")
)

// TodoService 待办事项服务
type TodoService struct {
    repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}

// Create 创建待办
func (s *TodoService) Create(userID uint, title, content string, startTime, dueTime *time.Time) (*model.Todo, error) {
    if title == "" {
        return nil, ErrTitleRequired
    }
    if content == "" {
        return nil, ErrContentRequired
    }

    todo := &model.Todo{
        UserID:    userID,
        Title:     title,
        Content:   content,
        View:      0,
        Status:    0, // TODO
        StartTime: startTime,
        DueTime:   dueTime,
    }

    if err := s.repo.Create(todo); err != nil {
        return nil, err
    }
    return todo, nil
}

