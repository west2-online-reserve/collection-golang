package service

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants"
)

func GetUserIDFromContext(c *app.RequestContext) int64 {
	return c.GetInt64(constants.ContextUserID)
}
