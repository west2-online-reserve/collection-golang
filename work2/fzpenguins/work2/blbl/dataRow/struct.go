package dataRow

type Body struct {
	Data Data        `json:"data"`
	Info interface{} `json:"info"`
}

type Data struct {
	Replies []Replies   `json:"replies"`
	Info    interface{} `json:"info"`
}

type Replies struct {
	Content Content     `json:"content"`
	Count   int         `json:"count"`
	Rpid    int64       `json:"rpid"`
	Mid     int64       `json:"mid"`
	Info    interface{} `json:"info"`
}

type Content struct {
	Message string      `json:"message"`
	Info    interface{} `json:"info"`
}
