package ctl

import (
	"net/http"
	"three/pkg/e"
	"three/types"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

type TaskItemList struct {
	Count int64                 `json:"count"`
	Page  int                   `json:"page"`
	Items []*types.TaskInfoResp `json:"items"`
}

func RespSuccess(code int) *Response {
	return &Response{
		Status: http.StatusOK,
		Data:   nil,
		Msg:    e.GetMsg(code),
		Error:  "",
	}
}

func RespSuccessWithData(data interface{}, code int) *Response {
	return &Response{
		Status: http.StatusOK,
		Data:   data,
		Msg:    e.GetMsg(code),
		Error:  "",
	}
}

func RespError(err error, code int) *Response {
	return &Response{
		Status: http.StatusOK,
		Data:   nil,
		Msg:    e.GetMsg(code),
		Error:  err.Error(),
	}
}

func RespErrorWithData(data interface{}, err error, code int) *Response {
	return &Response{
		Status: http.StatusOK,
		Data:   data,
		Msg:    e.GetMsg(code),
		Error:  err.Error(),
	}
}
