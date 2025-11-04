package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "memogo/biz/dal/model"
    redisClient "memogo/biz/dal/redis"

    "gorm.io/gorm"
)

// TodoRepository 待办事项数据访问层
type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

// 缓存键生成函数
func (r *TodoRepository) listCacheKey(userID uint, statusFilter string, page, pageSize int) string {
    return fmt.Sprintf("todos:list:user:%d:status:%s:page:%d:size:%d", userID, statusFilter, page, pageSize)
}

func (r *TodoRepository) searchCacheKey(userID uint, keyword string, page, pageSize int) string {
    return fmt.Sprintf("todos:search:user:%d:kw:%s:page:%d:size:%d", userID, keyword, page, pageSize)
}

func (r *TodoRepository) userCachePattern(userID uint) string {
    return fmt.Sprintf("todos:*:user:%d:*", userID)
}

// 缓存失效函数：删除某用户的所有待办缓存
func (r *TodoRepository) invalidateUserCache(userID uint) {
    if redisClient.RDB == nil {
        return
    }

    ctx := context.Background()
    pattern := r.userCachePattern(userID)

    // 扫描并删除匹配的键
    iter := redisClient.RDB.Scan(ctx, 0, pattern, 0).Iterator()
    for iter.Next(ctx) {
        redisClient.RDB.Del(ctx, iter.Val())
    }
}

// Create 新建待办
func (r *TodoRepository) Create(todo *model.Todo) error {
    if err := r.db.Create(todo).Error; err != nil {
        return err
    }
    // 清除该用户的缓存
    r.invalidateUserCache(todo.UserID)
    return nil
}

// UpdateStatusByID 按 ID 更新状态（限定用户）
func (r *TodoRepository) UpdateStatusByID(userID, id uint, status int32) (int64, error) {
    tx := r.db.Model(&model.Todo{}).
        Where("id = ? AND user_id = ?", id, userID).
        Update("status", status)
    if tx.Error != nil {
        return 0, tx.Error
    }
    // 清除该用户的缓存
    r.invalidateUserCache(userID)
    return tx.RowsAffected, nil
}

// UpdateAllStatus 按 from_status → to_status 批量更新（限定用户）
func (r *TodoRepository) UpdateAllStatus(userID uint, fromStatus, toStatus int32) (int64, error) {
    tx := r.db.Model(&model.Todo{}).
        Where("user_id = ? AND status = ?", userID, fromStatus).
        Update("status", toStatus)
    if tx.Error != nil {
        return 0, tx.Error
    }
    // 清除该用户的缓存
    r.invalidateUserCache(userID)
    return tx.RowsAffected, nil
}

// DeleteOne 删除单条（软删除，限定用户）
func (r *TodoRepository) DeleteOne(userID, id uint) (int64, error) {
    tx := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Todo{})
    if tx.Error != nil {
        return 0, tx.Error
    }
    // 清除该用户的缓存
    r.invalidateUserCache(userID)
    return tx.RowsAffected, nil
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
    if tx.Error != nil {
        return 0, tx.Error
    }
    // 清除该用户的缓存
    r.invalidateUserCache(userID)
    return tx.RowsAffected, nil
}

// ListTodos 分页查询（按状态筛选，可选）
func (r *TodoRepository) ListTodos(userID uint, statusFilter string, page, pageSize int) ([]model.Todo, int64, error) {
    var (
        todos []model.Todo
        total int64
    )

    // 尝试从缓存获取
    if redisClient.RDB != nil {
        cacheKey := r.listCacheKey(userID, statusFilter, page, pageSize)
        ctx := context.Background()

        cachedData, err := redisClient.RDB.Get(ctx, cacheKey).Result()
        if err == nil {
            // 缓存命中，解析数据
            type CachedResult struct {
                Todos []model.Todo `json:"todos"`
                Total int64        `json:"total"`
            }
            var cached CachedResult
            if json.Unmarshal([]byte(cachedData), &cached) == nil {
                return cached.Todos, cached.Total, nil
            }
        }
    }

    // 缓存未命中，查询数据库
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
    // 按创建时间升序：最早的备忘录显示在前面
    if err := q.Order("created_at ASC, id ASC").Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
        return nil, 0, err
    }

    // 将结果写入缓存（5分钟过期）
    if redisClient.RDB != nil {
        cacheKey := r.listCacheKey(userID, statusFilter, page, pageSize)
        ctx := context.Background()

        type CachedResult struct {
            Todos []model.Todo `json:"todos"`
            Total int64        `json:"total"`
        }
        cached := CachedResult{Todos: todos, Total: total}
        if data, err := json.Marshal(cached); err == nil {
            redisClient.RDB.Set(ctx, cacheKey, data, 5*time.Minute)
        }
    }

    return todos, total, nil
}

// SearchTodos 分页关键词查询（title/content 模糊匹配）
func (r *TodoRepository) SearchTodos(userID uint, keyword string, page, pageSize int) ([]model.Todo, int64, error) {
    var (
        todos []model.Todo
        total int64
    )

    // 尝试从缓存获取
    if redisClient.RDB != nil {
        cacheKey := r.searchCacheKey(userID, keyword, page, pageSize)
        ctx := context.Background()

        cachedData, err := redisClient.RDB.Get(ctx, cacheKey).Result()
        if err == nil {
            // 缓存命中，解析数据
            type CachedResult struct {
                Todos []model.Todo `json:"todos"`
                Total int64        `json:"total"`
            }
            var cached CachedResult
            if json.Unmarshal([]byte(cachedData), &cached) == nil {
                return cached.Todos, cached.Total, nil
            }
        }
    }

    // 缓存未命中，查询数据库
    q := r.db.Model(&model.Todo{}).
        Where("user_id = ?", userID).
        Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if page < 1 { page = 1 }
    if pageSize <= 0 { pageSize = 10 }
    offset := (page - 1) * pageSize
    // 按创建时间升序：最早的备忘录显示在前面
    if err := q.Order("created_at ASC, id ASC").Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
        return nil, 0, err
    }

    // 将结果写入缓存（5分钟过期）
    if redisClient.RDB != nil {
        cacheKey := r.searchCacheKey(userID, keyword, page, pageSize)
        ctx := context.Background()

        type CachedResult struct {
            Todos []model.Todo `json:"todos"`
            Total int64        `json:"total"`
        }
        cached := CachedResult{Todos: todos, Total: total}
        if data, err := json.Marshal(cached); err == nil {
            redisClient.RDB.Set(ctx, cacheKey, data, 5*time.Minute)
        }
    }

    return todos, total, nil
}

// ListTodosCursor 游标分页查询（用于高效遍历全部数据）
// cursor: 上一页最后一条的 ID，首次查询传 0
// 返回: todos列表, 下一页的cursor(0表示无下一页), hasMore(是否有更多数据), error
func (r *TodoRepository) ListTodosCursor(userID uint, statusFilter string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    var todos []model.Todo

    // 构建基础查询
    q := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)

    // 状态过滤
    switch statusFilter {
    case "done":
        q = q.Where("status = ?", 1)
    case "todo":
        q = q.Where("status = ?", 0)
    }

    // 游标过滤：因为是升序（旧→新），需要找比 cursor 更大的 ID
    if cursor > 0 {
        q = q.Where("id > ?", cursor)
    }

    // 查询 limit+1 条，用于判断是否还有下一页
    if err := q.Order("created_at ASC, id ASC").Limit(limit + 1).Find(&todos).Error; err != nil {
        return nil, 0, false, err
    }

    // 判断是否有更多数据
    hasMore := len(todos) > limit
    var nextCursor uint

    if hasMore {
        // 有下一页：返回 limit 条数据，nextCursor 为最后一条的 ID
        nextCursor = uint(todos[limit-1].ID)
        todos = todos[:limit]
    } else {
        // 没有下一页
        nextCursor = 0
    }

    return todos, nextCursor, hasMore, nil
}

// SearchTodosCursor 关键词游标分页查询
func (r *TodoRepository) SearchTodosCursor(userID uint, keyword string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    var todos []model.Todo

    // 构建查询
    q := r.db.Model(&model.Todo{}).
        Where("user_id = ?", userID).
        Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

    // 游标过滤
    if cursor > 0 {
        q = q.Where("id > ?", cursor)
    }

    // 查询 limit+1 条
    if err := q.Order("created_at ASC, id ASC").Limit(limit + 1).Find(&todos).Error; err != nil {
        return nil, 0, false, err
    }

    // 判断是否有更多数据
    hasMore := len(todos) > limit
    var nextCursor uint

    if hasMore {
        nextCursor = uint(todos[limit-1].ID)
        todos = todos[:limit]
    } else {
        nextCursor = 0
    }

    return todos, nextCursor, hasMore, nil
}
