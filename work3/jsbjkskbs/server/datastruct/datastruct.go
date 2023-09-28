package datastruct

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type TodolistBindJSONReceive struct {
	Title 	 string `json:"title"`
	Owner	 string `json:"owner"`
	Text     string `json:"todo"`
	Deadline string `json:"deadlline"`
	Status   bool   `json:"isdone"`
}

type TodolistBindJSONSend struct{
	Id		 int64	`json:"msgId"`
	Title 	 string `json:"title"`
	Owner	 string `json:"owner"`
	Text     string `json:"todo"`
	Addtime  string `json:"addtime"`
	Deadline string `json:"deadline"`
	Status   bool   `json:"isdone"`
}

type TodolistBindJSONSendArray struct{
	Todolist struct{
		List []TodolistBindJSONSend
	}`json:"todolist"`
}

type TodolistBindMysqlInsert struct{
	JsonData TodolistBindJSONReceive 
	Addtime  string
	Id		 int64
}

type TodolistBindMysqlSearch struct{
	Status   		bool 		`json:"isdone"` 
	Key		 		string		`json:"keyword"`
	SearchMethod	int 		`json:"method"`
	Username string
}
