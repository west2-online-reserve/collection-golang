package model

import (
	"Todov3/Redis"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type Task struct {
	gorm.Model
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index;not null"`
	Status    int    `gorm:"default:0"`
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64 `gorm:"default:0"`
}

func TaskViewKey(id uint) string {
	return fmt.Sprintf("view:task:%s", strconv.Itoa(int(id)))
}

func (task *Task) View() uint64 {
	c := context.Background()
	countStr, _ := Redis.RedisClient.Get(c, TaskViewKey(task.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

func (task *Task) AddView() {
	ctx := context.Background()
	Redis.RedisClient.Incr(ctx, TaskViewKey(task.ID))
}
