package datastruct

//在mysql版本不加除json以外的tag能正常转化interface{}到datastruct.User
//而在redis版本却要加与loginstruct/registerstruct一样的tag,应该不是redis的bug,是否为语言的问题?
type User struct {
	Username string `form:"username" json:"username" vd:"(len($)>0&&len($)<64); msg:'Illegal Username'"`
	Password string `form:"password" json:"password" vd:"(len($)>5&&len($)<16); msg:'Illegal Password'"`
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
