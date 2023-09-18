package model

import (
	"gorm.io/gorm"
	"strconv"
	"three/repository/cache"
)

type Task struct {
	gorm.Model
	Uid       uint
	Title     string `gorm:"index;not null"`
	Content   string `gorm:"type:longtext"`
	Status    int    `gorm:"default:0"`
	StartTime int64
	EndTime   int64 `gorm:"default:0"`
}

func (*Task) TableName() string {
	return "task"
}

func (t *Task) View() int {
	// 浏览次数
	_ = cache.RedisClient.SetNX(cache.TaskViewKey(t.ID), 1, 0)
	countStr, _ := cache.RedisClient.Get(cache.TaskViewKey(t.ID)).Result()
	count, _ := strconv.Atoi(countStr)
	return count
}

func (t *Task) AddView() {
	// 增加浏览次数
	cache.RedisClient.Incr(cache.TaskViewKey(t.ID))
	cache.RedisClient.ZIncrBy(cache.RankKey, 1, cache.TaskViewKey(t.ID))
}
