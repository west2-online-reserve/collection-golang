package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/west2-online-reserve/collection-golang/work3/biz/pack"
	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants"
	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants/jwt"
)

func AuthToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader(constants.AuthHeader))

		hlog.CtxInfof(ctx, "token: %+v clientIP: %v\n", token, c.ClientIP())

		claims, err := jwt.CheckToken(token)

		if err != nil {
			hlog.CtxErrorf(ctx, "check token err: %+v\n", err)
			pack.SendFailResponse(c, err)
			c.AbortWithStatus(401)
			return
		}

		hlog.CtxInfof(ctx, "claims: %+v\n", claims)

		token, err = jwt.CreateToken(claims.UserID)

		if err != nil {
			hlog.CtxInfof(ctx, "create token failed, client IP: %+v\n", c.ClientIP())
			pack.SendFailResponse(c, err)
			c.AbortWithStatus(401)
			return
		}

		c.Header(constants.AuthHeader, token)
		c.Set(constants.ContextUserID, claims.UserID)

		c.Next(ctx)
	}
}
