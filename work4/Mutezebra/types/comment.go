package types

type VideoCommentReq struct {
	VideoID  uint   `json:"video_id" form:"video_id"`
	Content  string `json:"content" form:"content"`
	ReplyID  uint   `json:"reply_id" form:"reply_id"`
	ReplyUid uint   `json:"reply_uid" form:"reply_uid,default:0"`
}
