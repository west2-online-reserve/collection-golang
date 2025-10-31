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

// UpdateStatusByID 按 ID 更新状态（限定用户）
func (r *TodoRepository) UpdateStatusByID(userID, id uint, status int32) (int64, error) {
    tx := r.db.Model(&model.Todo{}).
        Where("id = ? AND user_id = ?", id, userID).
        Update("status", status)
    return tx.RowsAffected, tx.Error
}

// UpdateAllStatus 按 from_status → to_status 批量更新（限定用户）
func (r *TodoRepository) UpdateAllStatus(userID uint, fromStatus, toStatus int32) (int64, error) {
    tx := r.db.Model(&model.Todo{}).
        Where("user_id = ? AND status = ?", userID, fromStatus).
        Update("status", toStatus)
    return tx.RowsAffected, tx.Error
}

// DeleteOne 删除单条（软删除，限定用户）
func (r *TodoRepository) DeleteOne(userID, id uint) (int64, error) {
    tx := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Todo{})
    return tx.RowsAffected, tx.Error
}

// DeleteByScope 按范围删除（done/todo/all，软删除，限定用户）
func (r *TodoRepository) DeleteByScope(userID uint, scope string) (int64, error) {
    q := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)
    switch scope {
    case "done":
        q = q.Where("status = ?", 1)
    case "todo":
        q = q.Where("status = ?", 0)
    case "all":
        // no extra filter
    default:
        // 未知 scope 交给上层校验，这里默认不执行
        return 0, nil
    }
    tx := q.Delete(&model.Todo{})
    return tx.RowsAffected, tx.Error
}

// ListTodos 分页查询（按状态筛选，可选）
func (r *TodoRepository) ListTodos(userID uint, statusFilter string, page, pageSize int) ([]model.Todo, int64, error) {
    var (
        todos []model.Todo
        total int64
    )
    q := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)
    switch statusFilter {
    case "done":
        q = q.Where("status = ?", 1)
    case "todo":
        q = q.Where("status = ?", 0)
    }
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if page < 1 { page = 1 }
    if pageSize <= 0 { pageSize = 10 }
    offset := (page - 1) * pageSize
    if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    return todos, total, nil
}

// SearchTodos 分页关键词查询（title/content 模糊匹配）
func (r *TodoRepository) SearchTodos(userID uint, keyword string, page, pageSize int) ([]model.Todo, int64, error) {
    var (
        todos []model.Todo
        total int64
    )
    q := r.db.Model(&model.Todo{}).
        Where("user_id = ?", userID).
        Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if page < 1 { page = 1 }
    if pageSize <= 0 { pageSize = 10 }
    offset := (page - 1) * pageSize
    if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
        return nil, 0, err
    }
    return todos, total, nil
}
