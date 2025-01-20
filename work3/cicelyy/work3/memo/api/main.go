package api

import (
	"encoding/json"
	"fmt"
	"memo/serializer" 

	// "github.com/go-playground/validator/v10"
	// "gopkg.in/go-playgrond/validator.v8"
)

//处理错误并返回一个标准化的错误响
func ErrorResponse(err error) serializer.Response{
	//先检查类型是否匹配
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 40001,
			Msg: "Json类型不匹配",
			Error: fmt.Sprint(err),
		}
	}
	return serializer.Response{
        Status: 40001,
        Msg:    "参数错误",
        Error:  fmt.Sprint(err),
    }
}