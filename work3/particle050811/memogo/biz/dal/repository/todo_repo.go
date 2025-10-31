package repository

import (
    "memogo/biz/dal/model"
    "gorm.io/gorm"
)

// TodoRepository 待办事项数据访问层
type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

// Create 新建待办
func (r *TodoRepository) Create(todo *model.Todo) error {
    return r.db.Create(todo).Error
}

