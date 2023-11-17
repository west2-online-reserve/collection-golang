package datastruct

//在mysql版本不加除json以外的tag能正常转化interface{}到datastruct.User
//而在redis版本却要加与loginstruct/registerstruct一样的tag,应该不是redis的bug,是否为语言的问题?
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ShortResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SendingJSONFormat struct {
	Status       int             `json:"status"`
	Page         int             `json:"page"`
	ItemsPerPage int             `json:"itemsPerPage"`
	Isend        bool            `json:"isend"`
	Data         SendingJSONData `json:"data"`
	Message      string          `json:"message"`
	Error        string          `json:"error"`
}

type SendingJSONData struct {
	Items []TodolistBindJSONSend `json:"items"`
	Count int                    `json:"count"`
	Total int                    `json:"total"`
}

type TodolistBindJSONReceive struct {
	Title    string `json:"title" example:"标题"`
	Owner    string `json:"owner" swaggerignore:"true"`
	Text     string `json:"text" example:"文本"`
	Deadline string `json:"deadline" example:"2077-01-01 01:01:01"`
	Status   bool   `json:"isdone" swaggerignore:"true"`
}

type TodolistBindJSONSend struct {
	Id       int64  `json:"msgId"`
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	Text     string `json:"text"`
	Addtime  string `json:"addtime"`
	Deadline string `json:"deadline"`
	Status   bool   `json:"isdone"`
}

type TodolistBindJSONSendArray struct {
	Todolist struct {
		List []TodolistBindJSONSend
	} `json:"todolist"`
}

type TodolistBindRedisCondition struct {
	Keyword string  `json:"keyword" example:"我超OP"`
	Status  bool    `json:"isdone" example:"false"`
	Idlist  []int64 `json:"idlist"`
	Method  int     `json:"method" example:"1"`
}

type TodolistBindRedisUpdate struct {
	Idlist    []int64 `json:"idlist"`
	NewStatus bool    `json:"isdone" example:"false"`
	Method    int     `json:"method" example:"1"`
}

type ProcessStruct struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	Text     string `json:"text"`
	Addtime  string `json:"addtime"`
	Deadline string `json:"deadline"`
	Isdone   bool   `json:"isdone"`
}