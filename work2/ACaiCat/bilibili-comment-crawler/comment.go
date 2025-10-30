package main

import "time"

type Comment struct {
	Index           int `gorm:"primaryKey;index;autoincrement"`
	AuthorName      string
	AuthorUID       string
	Time            time.Time
	Like            int
	Content         string
	IsReply         bool
	ParentCommentID int
	Relies          []Comment `gorm:"foreignKey:ParentCommentID"`
}
