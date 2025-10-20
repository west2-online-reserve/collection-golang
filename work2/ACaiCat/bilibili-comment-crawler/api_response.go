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
	Page    Page    `json:"page"`    //页码
	Replies []Reply `json:"replies"` //评论
}

// Page 页码信息
type Page struct {
	Num    int `json:"num"`    //当前页码
	Size   int `json:"size"`   //当前页评论数
	Count  int `json:"count"`  //展示评论数
	Acount int `json:"acount"` //总评论数
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
