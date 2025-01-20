//定义了一个任务的数据库模型，包含了任务的基本信息
package model

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model //内嵌的结构体
	User User `gorm:"ForeignKey:Uid"` //用户
	Uid uint `gorm:"not null"` //ID
	Title string `gorm:"index:not null"`//有索引
	Status int `gorm:"default:'0'"` //状态，0是未完成
	Content string `gorm:"type:longtext"` //内容
	StartTime int64 //备忘录开始时间 时间戳
	EndTime int64 //备忘录完成时间 时间戳 
} 