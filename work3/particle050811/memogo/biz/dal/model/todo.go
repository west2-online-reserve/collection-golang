package model

import (
    "time"

    "gorm.io/gorm"
)

// Todo 待办事项模型
type Todo struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    UserID  uint   `gorm:"index;not null" json:"user_id"`
    Title   string `gorm:"not null;size:200" json:"title"`
    Content string `gorm:"type:text;not null" json:"content"`

    View   int32 `gorm:"not null;default:0" json:"view"`
    Status int32 `gorm:"not null;default:0" json:"status"` // 0: TODO, 1: DONE

    StartTime *time.Time `json:"start_time"`
    EndTime   *time.Time `json:"end_time"`
    DueTime   *time.Time `json:"due_time"`
}

func (Todo) TableName() string { return "todos" }

