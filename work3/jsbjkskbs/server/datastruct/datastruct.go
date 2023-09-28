package datastruct

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ShortResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SendingJSONFormat struct {
	Status  int             `json:"status"`
	Data    SendingJSONData `json:"data"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}

type SendingJSONData struct {
	Items      []TodolistBindJSONSend `json:"items"`
	TotalItems int                    `json:"count"`
}

type TodolistBindJSONReceive struct {
	Title    string `json:"title" example:"标题"`
	Owner    string `json:"owner" swaggerignore:"true"`
	Text     string `json:"todo" example:"文本"`
	Deadline string `json:"deadline" example:"2077-01-01 01:01:01"`
	Status   bool   `json:"isdone" swaggerignore:"true"`
}

type TodolistBindJSONSend struct {
	Id       int64  `json:"msgId"`
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	Text     string `json:"todo"`
	Addtime  string `json:"addtime"`
	Deadline string `json:"deadline"`
	Status   bool   `json:"isdone"`
}

type TodolistBindJSONSendArray struct {
	Todolist struct {
		List []TodolistBindJSONSend
	} `json:"todolist"`
}

type TodolistBindMysqlInsert struct {
	JsonData TodolistBindJSONReceive
	Addtime  string
	Id       int64
}

type TodolistBindMysqlSearch struct {
	Status       bool   `json:"isdone" example:"true"`
	Key          string `json:"keyword" example:"OP"`
	SearchMethod int    `json:"method" example:"1"`
	Username     string `swaggerignore:"true"`
}

type TodolistBindMysqlDelete struct {
	Status       bool    `json:"isdone" example:"true"`
	Idlist       []int64 `json:"idlist" example:[]`
	DeleteMethod int     `json:"method" example:"1"`
	Username     string  `swaggerignore:"true"`
}

type TodolistBindMysqlModify struct {
	Status       bool    `json:"isdone" example:"true"`
	Idlist       []int64 `json:"idlist" example:[]`
	ModifyMethod int     `json:"method" example:"1"`
	Username     string  `swaggerignore:"true"`
}
