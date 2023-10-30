package model

import (
	"fmt"
	"four/consts"
	"gorm.io/gorm"
)

// Comment 一百个视频一张表
type Comment struct {
	gorm.Model
	VideoID  uint
	Root     uint // 根评论id
	ReplyID  uint // 如果为0说明是对视频的评论，否则是对其他评论的回复
	Uid      uint
	ReplyUid uint   // 被回复的人的uid
	Content  string `gorm:"type:text"`
}

// TableExist 判断表是否存在
func (c *Comment) TableExist() bool {
	result := c.VideoID % consts.EachVideoRecordACommentTable
	if result == 1 {
		return false
	}
	return true
}

func (c *Comment) InsertNewCommentSQL(tableName string) string {
	return fmt.Sprintf("%s%s%s,%d,%d,%d,'%s',%d,%d);", consts.InsertNewComment1, tableName, consts.InsertNewComment2, c.VideoID, c.ReplyID, c.Uid, c.Content, c.ReplyUid, c.Root)
}

func (c *Comment) FindCommentRootSQL(tableName string) string {
	return fmt.Sprintf(fmt.Sprintf("SELECT root FROM %s WHERE id=%d LIMIT 1", tableName, c.ReplyID))
}
