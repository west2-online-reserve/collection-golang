package response

import "net/http"

type TaskResp struct {
	ID        uint   `json:"id" example:"1"`       // 任务ID
	Title     string `json:"title" example:"吃饭"`   // 题目
	Content   string `json:"content" example:"睡觉"` // 内容
	View      uint64 `json:"view" example:"32"`    // 浏览量
	Status    int    `json:"status" example:"0"`   // 状态(0未完成，1已完成)
	CreatedAt int64  `json:"created_at"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

func TaskListRep(items interface{}, total int64) Response {
	return Response{
		Status: http.StatusOK,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Message: "ok",
		Error:   "",
	}
}
