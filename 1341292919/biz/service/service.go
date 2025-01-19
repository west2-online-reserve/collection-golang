package service

import (
	"Demo/pkg/constants"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetUserIDFromContext(c *app.RequestContext) int64 {
	return c.GetInt64(constants.ContextUserId)
}
