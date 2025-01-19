package middleware

import (
	"Demo/biz/pack"
	"Demo/pkg/constants"
	"Demo/pkg/constants/jwt"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AutoToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader(constants.AuthHeader))
		hlog.CtxInfof(ctx, "token: %+v clientIP : %v\n", token, c.ClientIP())

		claims, err := jwt.CheckToken(token)
		if err != nil {
			hlog.CtxErrorf(ctx, "check token err: %+v\n", err)
			pack.SendFailResponse(c, err)
			c.AbortWithStatus(401)
			return
		}
		hlog.CtxInfof(ctx, "token: %+v\n", claims)

		token, err = jwt.CreateToken(claims.UserID)

		if err != nil {
			hlog.CtxErrorf(ctx, "create token failed client IP: %+v\n", c.ClientIP())
			pack.SendFailResponse(c, err)
			c.AbortWithStatus(401)
			return
		}

		c.Header(constants.AuthHeader, token)
		c.Set(constants.ContextUserId, claims.UserID)

		c.Next(ctx)
	}
}
