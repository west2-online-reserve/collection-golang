package api

import (
	"encoding/json"
	"three/pkg/ctl"
	"three/pkg/e"
)

func ErrorResponse(err error) *ctl.Response {
	_, ok := err.(*json.UnmarshalTypeError)
	if ok {
		return ctl.RespError(err, e.JsonUnmarshalFailed)
	}
	return ctl.RespError(err, e.ERROR)
}
