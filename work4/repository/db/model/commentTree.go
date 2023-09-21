package model

import "gorm.io/gorm"

type CommentTree struct {
	gorm.Model
	CommentId uint
	ReplyId   uint
}
