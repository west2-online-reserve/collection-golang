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
	TotalItems int                 `json:"count"`
}

type TodolistBindJSONReceive struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	Text     string `json:"todo"`
	Deadline string `json:"deadline"`
	Status   bool   `json:"isdone"`
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
	Status       bool   `json:"isdone"`
	Key          string `json:"keyword"`
	SearchMethod int    `json:"method"`
	Username     string
}

type TodolistBindMysqlDelete struct {
	Status       bool    `json:"isdone"`
	Idlist       []int64 `json:"idlist"`
	DeleteMethod int     `json:"method"`
	Username     string
}

type TodolistBindMysqlModify struct {
	Status       bool    `json:"isdone"`
	Idlist       []int64 `json:"idlist"`
	ModifyMethod int     `json:"method"`
	Username     string
}
