package middleware

import (
	"context"
	"errors"
	"four/consts"
	"four/pkg/ctl"
	"four/pkg/e"
	"four/pkg/myutils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var code int
		code = e.SUCCESS
		accessToken := string(c.GetHeader("access_token"))
		refreshToken := string(c.GetHeader("refresh_token"))
		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(http.StatusBadRequest, utils.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  errors.New("token is empty").Error(),
			})
			c.Abort()
			return
		}

		newAToken, newRToken, err := myutils.CheckToken(accessToken, refreshToken)

		if err != nil {
			code = e.CheckTokenFailed
			c.JSON(http.StatusInternalServerError, utils.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		claims, err, _ := myutils.ParseToken(newAToken)
		if err != nil {
			code = e.ParseTokenFailed
			c.JSON(http.StatusInternalServerError, utils.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		setHeader(c, newAToken, newRToken)
		ctx = ctl.NewContext(ctx, &ctl.UserInfo{ID: claims.ID, UserName: claims.UserName})
		c.Next(ctx)
	}
}
func setHeader(c *app.RequestContext, aToken, rToken string) {
	c.Header(consts.AccessToken, aToken)
	c.Header(consts.RefreshToken, rToken)
}
