package serializer

//基础序列化器 通用基础结构体
type Response struct {
	Status int `json:"status"` //状态
	Data interface{} `json:"data"` 
	Msg string `json:"msg"` //返回的信息
	Error string `json:"error"` //返回的错误
}

//带有token的Data数据结构
type TokenData struct {
	User interface{} `json:"user"`
	Token string `json:"token"`
}

//DataList 带有总数的Data结构
type DataList struct {
	Item interface{} `json:"item"`
	Total uint  `json:"total"`
}

//带总数的返回
func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item: items,
			Total: total,
		},
		Msg: "ok",
	}
}