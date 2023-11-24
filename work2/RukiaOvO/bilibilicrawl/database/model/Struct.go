package model

type DbComments struct {
	Num         int64  `gorm:"autoIncrement; primaryKey"`
	MainUname   string `gorm:"not null"`
	MainContent string `gorm:"not null"`
	RepUname1   string `gorm:"default null"`
	RepContent1 string `gorm:"default null"`
	RepUname2   string `gorm:"default null"`
	RepContent2 string `gorm:"default null"`
	RepUname3   string `gorm:"default null"`
	RepContent3 string `gorm:"default null"`
}

type BilibiliComment struct {
	Data Data `json:"data"` //评论数据
}
type Data struct {
	MainReply []MainReply `json:"replies"` //主评论
}
type MainReply struct {
	Content Content `json:"content"` //主评论内容
	Member  Member  `json:"member"`  //发表主评论的用户
	Replies []Reply `json:"replies"` //回复主评论的次评论
}
type Reply struct {
	Content Content `json:"content"` //评论文本
	Member  Member  `json:"member"`  //发表次评论的用户
}
type Member struct {
	Uname string `json:"uname"` //用户ID
}
type Content struct {
	Message string `json:"message"` //评论的内容
}
