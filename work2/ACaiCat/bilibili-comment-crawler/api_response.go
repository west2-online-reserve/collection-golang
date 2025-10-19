package main

// APIResponse API响应
type APIResponse struct {
	Code    int    `json:"code"`    //状态码
	Message string `json:"message"` //错误信息
	TTL     int    `json:"ttl"`     //缓存时间
	Data    Data   `json:"data"`    //负载数据
}

// Data 负载数据
type Data struct {
	Cursor  Cursor  `json:"cursor"`  //游标
	Replies []Reply `json:"replies"` //评论
}

// PaginationReply 页码信息
type PaginationReply struct {
	NextOffset string `json:"next_offset"` // 下一页评论的偏移量，is_end为true时不存在
}
type Cursor struct {
	IsBegin         bool            `json:"is_begin"`         //第一页
	IsEnd           bool            `json:"is_end"`           //最后一页
	PaginationReply PaginationReply `json:"pagination_reply"` //页码信息
	AllCount        int             `json:"all_count"`        //评论总数
	Name            string          `json:"name"`             //评论类型 最新评论|最热评论
}

// Content 内容
type Content struct {
	Message string `json:"message"` //纯文本内容
}

// Reply 评论|回复
type Reply struct {
	Oid     int     `json:"oid"`               //可能是视频的唯一ID
	Like    int     `json:"like"`              //点赞数
	Ctime   int64   `json:"ctime"`             //发布时间戳
	Content Content `json:"content,omitempty"` //内容
	Member  Member  `json:"member"`            //发送者信息
	Replies []Reply `json:"replies"`           //回复
}

// Member 评论者
type Member struct {
	Mid   string `json:"mid"`   //UID
	Uname string `json:"uname"` //昵称
}
