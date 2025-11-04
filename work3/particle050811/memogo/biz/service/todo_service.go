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

// UpdateTodoStatus 更新单条状态
func (s *TodoService) UpdateTodoStatus(userID, id uint, status int32) (int64, error) {
    return s.repo.UpdateStatusByID(userID, id, status)
}

// UpdateAllStatus 批量更新状态 from → to
func (s *TodoService) UpdateAllStatus(userID uint, fromStatus, toStatus int32) (int64, error) {
    return s.repo.UpdateAllStatus(userID, fromStatus, toStatus)
}

// DeleteOne 删除单条
func (s *TodoService) DeleteOne(userID, id uint) (int64, error) {
    return s.repo.DeleteOne(userID, id)
}

// DeleteByScope 范围删除
func (s *TodoService) DeleteByScope(userID uint, scope string) (int64, error) {
    return s.repo.DeleteByScope(userID, scope)
}

// ListTodos 分页查询（status: "todo"|"done"|"all"|""）
func (s *TodoService) ListTodos(userID uint, status string, page, pageSize int) ([]model.Todo, int64, error) {
    // 统一页大小限制
    if page < 1 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    } else if pageSize > 50 {
        pageSize = 50
    }
    return s.repo.ListTodos(userID, status, page, pageSize)
}

// SearchTodos 关键词分页查询
func (s *TodoService) SearchTodos(userID uint, keyword string, page, pageSize int) ([]model.Todo, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    } else if pageSize > 50 {
        pageSize = 50
    }
    return s.repo.SearchTodos(userID, keyword, page, pageSize)
}

// ListTodosCursor 游标分页查询（用于高效遍历全部数据，O(n) 复杂度）
func (s *TodoService) ListTodosCursor(userID uint, status string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    // 限制每次查询的最大数量
    if limit <= 0 {
        limit = 10
    } else if limit > 100 {
        limit = 100 // 游标分页可以允许更大的 limit
    }
    return s.repo.ListTodosCursor(userID, status, cursor, limit)
}

// SearchTodosCursor 关键词游标分页查询
func (s *TodoService) SearchTodosCursor(userID uint, keyword string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    if limit <= 0 {
        limit = 10
    } else if limit > 100 {
        limit = 100
    }
    return s.repo.SearchTodosCursor(userID, keyword, cursor, limit)
}
