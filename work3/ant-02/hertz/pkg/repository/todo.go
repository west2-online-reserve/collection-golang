package repository

import (
	"encoding/json"
	"hertz/database"
	"hertz/pkg/model"
	"log"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

const (
	baseKey string = "todo:"
)

type todoRepository struct {
	db *gorm.DB
	r  database.RedisClient
	mu sync.RWMutex
}

type TodoRepository interface {
	Create(todo *model.Todo) bool
	UpdateStatusById(uid uint64, id uint64, status uint16) bool
	UpdateFinishedStatus(uid uint64) bool
	UpdateNotFinishedStatus(uid uint64) bool
	QueryByStatus(uid uint64, status uint16, pageNum, pageSize int) ([]*model.Todo, int64)
	QueryByKey(uid uint64, key string, pageNum, pageSize int) ([]*model.Todo, int64)
	Query(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64)
	DeleteById(id uint64, uid uint64) bool
	DeleteByStatus(uid uint64, status uint16) bool
	DeleteAll(uid uint64) bool
}

func NewTodoRepository(db *gorm.DB, redisClient database.RedisClient) TodoRepository {
	return &todoRepository{
		db: db,
		r:  redisClient,
	}
}

// 缓存键生成函数
func (tr *todoRepository) cacheKeys(uid uint64) []string {
	return []string{
		baseKey + "all:" + strconv.FormatUint(uid, 10),      // 全部todo列表
		baseKey + "status:0:" + strconv.FormatUint(uid, 10), // 状态0的列表
		baseKey + "status:1:" + strconv.FormatUint(uid, 10), // 状态1的列表
	}
}

// 清除用户相关的所有缓存
func (tr *todoRepository) clearUserCache(uid uint64) bool {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	keys := tr.cacheKeys(uid)
	if err := tr.r.Del(keys); err != nil {
		log.Printf("删除缓存失败: %v", err)
		return false
	}
	return true
}

// 从Redis List获取分页数据（修复切片传递问题）
func (tr *todoRepository) queryTodoRange(key string, first, second int64) ([]*model.Todo, int64, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	// 检查key是否存在
	exists, err := tr.r.Exists(key)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, nil
	}

	// 获取列表长度
	total, err := tr.r.LLen(key)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	todoJson, err := tr.r.LRange(key, first, second)
	if err != nil {
		return nil, 0, err
	}

	// 解析数据
	todos := make([]*model.Todo, 0, len(todoJson))
	for _, s := range todoJson {
		var todo model.Todo
		if err := json.Unmarshal([]byte(s), &todo); err != nil {
			return nil, 0, err
		}
		todos = append(todos, &todo)
	}

	return todos, total, nil
}

// 设置Todo列表到缓存
func (tr *todoRepository) setTodoList(key string, todos []*model.Todo) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	// 先删除旧的列表
	if err := tr.r.Del([]string{key}); err != nil {
		return err
	}

	// 批量插入新数据
	for _, todo := range todos {
		s, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		if err := tr.r.RPush(key, s); err != nil {
			return err
		}
	}

	return nil
}

func (tr *todoRepository) Create(todo *model.Todo) bool {
	// 先清除缓存
	if !tr.clearUserCache(todo.Uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	// 创建数据库记录
	if err := tr.db.Create(todo).Error; err != nil {
		log.Printf("创建代办事项失败: %v", err)
		return false
	}

	return true
}

func (tr *todoRepository) UpdateStatusById(uid uint64, id uint64, status uint16) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Model(&model.Todo{}).
		Where("uid = ?", uid).
		Where("id = ?", id).
		Where("end_time > ?", time.Now()).
		Update("status", status).Error; err != nil {
		log.Printf("修改状态失败: %v", err)
		return false
	}
	return true
}

func (tr *todoRepository) UpdateFinishedStatus(uid uint64) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Model(&model.Todo{}).
		Where("uid = ?", uid).
		Where("end_time > ?", time.Now()).
		Where("status = 0").
		Update("status", 1).Error; err != nil {
		log.Printf("修改状态失败: %v", err)
		return false
	}
	return true
}

func (tr *todoRepository) UpdateNotFinishedStatus(uid uint64) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Model(&model.Todo{}).
		Where("uid = ?", uid).
		Where("end_time > ?", time.Now()).
		Where("status = 1").
		Update("status", 0).Error; err != nil {
		log.Printf("修改状态失败: %v", err)
		return false
	}
	return true
}

func (tr *todoRepository) QueryByStatus(uid uint64, status uint16, pageNum, pageSize int) ([]*model.Todo, int64) {
	key := baseKey + "status:" + strconv.FormatUint(uint64(status), 10) + ":" + strconv.FormatUint(uid, 10)
	first, second := int64((pageNum-1)*pageSize), int64(pageNum*pageSize-1)

	// 尝试从缓存获取
	todos, total, err := tr.queryTodoRange(key, first, second)
	if err != nil {
		log.Printf("从缓存查询失败: %v", err)
		// 继续查询数据库
	}

	if todos != nil {
		return todos, total
	}

	// 缓存未命中，查询数据库
	var res []*model.Todo
	var count int64

	// 使用事务确保计数和查询的一致性
	err = tr.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&model.Todo{}).
			Where("uid = ?", uid).
			Where("status = ?", status)

		if err := query.Count(&count).Error; err != nil {
			return err
		}

		if err := query.Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&res).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("查询代办事项失败: %v", err)
		return nil, 0
	}

	// 异步更新缓存
	go func() {
		if err := tr.setTodoList(key, res); err != nil {
			log.Printf("更新缓存失败: %v", err)
		}
	}()

	return res, count
}

func (tr *todoRepository) QueryByKey(uid uint64, keyword string, pageNum, pageSize int) ([]*model.Todo, int64) {
	var res []*model.Todo
	var count int64

	// 关键词搜索不缓存，因为组合太多
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&model.Todo{}).
			Where("uid = ?", uid).
			Where("title LIKE ?", "%"+keyword+"%")

		if err := query.Count(&count).Error; err != nil {
			return err
		}

		if err := query.Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&res).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("查询代办事项失败: %v", err)
		return nil, 0
	}
	return res, count
}

func (tr *todoRepository) Query(uid uint64, pageNum, pageSize int) ([]*model.Todo, int64) {
	key := baseKey + "all:" + strconv.FormatUint(uid, 10)
	first, second := int64((pageNum-1)*pageSize), int64(pageNum*pageSize-1)

	// 尝试从缓存获取
	todos, total, err := tr.queryTodoRange(key, first, second)
	if err != nil {
		log.Printf("从缓存查询失败: %v", err)
		// 继续查询数据库
	}

	if todos != nil {
		return todos, total
	}

	// 缓存未命中，查询数据库
	var res []*model.Todo
	var count int64

	err = tr.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&model.Todo{}).
			Where("uid = ?", uid)

		if err := query.Count(&count).Error; err != nil {
			return err
		}

		if err := query.Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&res).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("查询代办事项失败: %v", err)
		return nil, 0
	}

	// 异步更新缓存
	go func() {
		if err := tr.setTodoList(key, res); err != nil {
			log.Printf("更新缓存失败: %v", err)
		}
	}()

	return res, count
}

func (tr *todoRepository) DeleteById(id uint64, uid uint64) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Where("uid = ?", uid).
		Delete(&model.Todo{}, id).Error; err != nil {
		log.Printf("按id删除代办事件失败: %v", err)
		return false
	}
	return true
}

func (tr *todoRepository) DeleteByStatus(uid uint64, status uint16) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Where("status = ?", status).
		Where("uid = ?", uid).
		Delete(&model.Todo{}).Error; err != nil {
		log.Printf("按状态删除代办事件失败: %v", err)
		return false
	}
	return true
}

func (tr *todoRepository) DeleteAll(uid uint64) bool {
	// 先清除缓存
	if !tr.clearUserCache(uid) {
		log.Printf("清除缓存失败，但继续数据库操作")
	}

	if err := tr.db.Where("uid = ?", uid).
		Delete(&model.Todo{}).Error; err != nil {
		log.Printf("删除全部代办事件失败: %v", err)
		return false
	}
	return true
}
