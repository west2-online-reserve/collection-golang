package types

// 增
type CreateTaskRequest struct {
	Title   string `json:"title" form:"title" binding:"required,max=100"`
	Content string `json:"content" form:"content" binding:"required,max=100"`
	Status  int    `json:"status" form:"status" default:"0"`
}

// 查
type ShowAllTaskRequest struct {
	Limit int `json:"limit" form:"limit"`
	Start int `json:"start" form:"start"`
}

// 输入要查询的条件：代办为0，已完成为1
type ShowAllTaskRequestWithCondition struct {
	Condition int `json:"condition" form:"condition"`
	Limit     int `json:"limit" form:"limit"`
	Start     int `json:"start" form:"start"`
}

type SearchTasksRequest struct {
	Info  string `json:"info" form:"info" binding:"required"`
	Limit int    `json:"limit" form:"limit"`
	Start int    `json:"start" form:"start"`
}

// 改
// 只修改status
type UpdateTaskRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type UpdateALlTasksRequest struct {
	Condition int `json:"condition" form:"condition"`
}

// 删
type DeleteTaskRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteTasksRequestWithCondition struct {
	Condition int `json:"condition" form:"condition"`
}
