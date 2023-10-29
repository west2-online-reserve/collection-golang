package response

import "net/http"

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
	Error   string      `json:"error"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total int64       `json:"total"`
}

func SuccessResponseWithData(data interface{}) Response {
	return Response{
		Status: http.StatusOK,
		Data: DataList{
			Item:  data,
			Total: 1,
		},
		Message: "操作成功",
		Error:   "",
	}
}

func SuccessResponse() Response {
	return Response{
		Status:  http.StatusOK,
		Data:    nil,
		Message: "ok",
		Error:   "",
	}
}

func BadResponse(msg string) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Data:    nil,
		Message: msg,
		Error:   "",
	}
}
