package ctl

import "four/pkg/e"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

func RespSuccess(code int) *Response {
	return &Response{
		Status: code,
		Data:   nil,
		Msg:    e.GetMsg(code),
		Error:  "",
	}
}

func RespSuccessWithData(code int, data interface{}) *Response {
	return &Response{
		Status: code,
		Data:   data,
		Msg:    e.GetMsg(code),
		Error:  "",
	}
}

func RespError(code int, err error) *Response {
	return &Response{
		Status: code,
		Data:   nil,
		Msg:    e.GetMsg(code),
		Error:  err.Error(),
	}
}

func RespErrorWithData(code int, err error, data interface{}) *Response {
	return &Response{
		Status: code,
		Data:   data,
		Msg:    e.GetMsg(code),
		Error:  err.Error(),
	}
}
