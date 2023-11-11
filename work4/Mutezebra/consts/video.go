package consts

import "time"

const (
	//EachVideoRecordACommentTable 每多少个视频的一个评论表
	EachVideoRecordACommentTable = 1000

	//EachUserRecordAFansTable 每多少个用户的一个粉丝表
	EachUserRecordAFansTable = 1000

	EveryPageSize = 2 * MB
	MaxVideoSize  = 50 * MB

	//RegularClean 定时清理内存的时间
	RegularClean = 15 * time.Second
)
