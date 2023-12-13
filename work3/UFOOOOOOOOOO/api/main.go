package api

import (
	// "encoding"
	"encoding/json"
	"fmt"
	// "todolist/conf"
	"todolist/serializer"

	// "github.com/go-playground/validator/v10"
	// "gopkg.in/go-playgrond/validator.v8"
)

func ErrorResponse(err error) serializer.Response{
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