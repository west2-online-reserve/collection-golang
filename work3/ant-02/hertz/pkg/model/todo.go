package model

import "time"

type Todo struct {
	ID        uint64    `json:"id"` // 待办事项ID
	Uid       uint64    `json:"uid"`
	Title     string    `json:"title"`      // 主题
	Content   string    `json:"content"`    // 内容
	View      uint64    `json:"view"`       // 访问次数
	Status    uint16    `json:"status"`     // 状态：1=进行中, 2=已完成, 3=已取消...
	CreatedAt time.Time `json:"created_at"` // Unix 时间戳（秒）
	StartTime time.Time `json:"start_time"` // Unix 时间戳（秒）
	EndTime   time.Time `json:"end_time"`   // Unix 时间戳（秒），0 表示未结束
}
