package types

type TaskCreateReq struct {
	Status  int    `json:"status" form:"status"` // 0 is live, 1 is unlived
	Content string `json:"content" form:"content"`
	Title   string `json:"title" form:"title"`
}

type TaskUpdateReq struct {
	Id      uint   `json:"id" form:"id"`
	Status  int    `json:"status" form:"status"` // 0 is live, 1 is unlived
	Content string `json:"content" form:"content"`
	Title   string `json:"title" form:"form"`
}

type TaskShowReq struct {
	Id uint `json:"id" form:"id"`
}

type TaskListReq struct {
	Limit int `json:"limit" form:"limit"`
	Start int `json:"start" form:"start"`
}

type TaskSearchReq struct {
	Text   string `json:"text" form:"text"`
	Status int    `json:"status" form:"status"`
	Start  int    `json:"start" form:"start"`
}

type TaskDeleteReq struct {
	Id     uint `json:"id" form:"id"`
	Status int  `json:"status" form:"status"`
}

type TaskInfoResp struct {
	Id        uint   `json:"id,omitempty"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	View      int    `json:"view,omitempty"`
	Status    int    `json:"status"`
	CreateAt  int64  `json:"create_at"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
