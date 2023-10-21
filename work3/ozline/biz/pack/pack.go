package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/west2-online-reserve/collection-golang/work3/biz/model/model"
)

type Base struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Response struct {
	Base Base `json:"base"`
}

func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}

func BuildBaseResp(err error) *model.BaseResp {
	if err == nil {
		return &model.BaseResp{
			Code: 10000,
			Msg:  "ok",
		}
	}

	return &model.BaseResp{
		Code: 10001,
		Msg:  err.Error(),
	}
}

func SendFailResponse(c *app.RequestContext, err error) {
	SendResponse(c, BuildBaseResp(err))
}
