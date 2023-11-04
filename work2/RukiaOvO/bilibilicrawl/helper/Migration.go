package helper

import (
	"bilibilicrawl/database/model"
	"gorm.io/gorm"
)

func DataMigration(y model.MainReply, db *gorm.DB) {
	dbData := model.DbComments{}

	dbData.MainUname = y.Member.Uname
	dbData.MainContent = y.Content.Message

	for len(y.Replies) < 3 {
		var temp model.Reply
		temp.Member.Uname = "null"
		temp.Content.Message = "null"
		y.Replies = append(y.Replies, temp)
	}

	dbData.RepUname1 = y.Replies[0].Member.Uname
	dbData.RepContent1 = y.Replies[0].Content.Message
	dbData.RepUname2 = y.Replies[1].Member.Uname
	dbData.RepContent2 = y.Replies[1].Content.Message
	dbData.RepUname3 = y.Replies[2].Member.Uname
	dbData.RepContent3 = y.Replies[2].Content.Message

	db.Create(&dbData)
}
