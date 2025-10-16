package main

import (
	"log"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CommentModel 是数据库存储模型
type CommentModel struct {
	Rpid      int64  `gorm:"primaryKey"` // 主键
	Uname     string `gorm:"size:128"`   // 用户名
	Message   string `gorm:"type:text"`  // 评论内容
	Ctime     int64  // 创建时间
	Likes     int    // 点赞数
	CreatedAt time.Time
	UpdatedAt time.Time
}

// openDB 初始化并返回 *gorm.DB
func openDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	if err := db.AutoMigrate(&CommentModel{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	log.Println("✅ GORM 数据库已初始化 comments.db")
	return db
}

// SaveComment 保存单条评论（存在则更新）
func SaveComment(db *gorm.DB, c Comment) {
	m := CommentModel{
		Rpid:    c.Rpid,
		Uname:   c.Member.Uname,
		Message: strings.ReplaceAll(c.Content.Message, "\n", " "),
		Ctime:   c.Ctime,
		Likes:   c.Like,
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rpid"}}, // 主键冲突时更新
		DoUpdates: clause.AssignmentColumns([]string{"uname", "message", "ctime", "likes", "updated_at"}),
	}).Create(&m).Error

	if err != nil {
		log.Printf("❌ 插入失败: %v", err)
	} else {
		log.Printf("💾 已保存评论 #%d 来自 %s", m.Rpid, m.Uname)
	}
}
