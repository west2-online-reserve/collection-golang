package api

import (
	"context"
	"four/pkg/log"
	"four/service"
	"four/types"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

func UserRegisterHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserRegisterReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Register(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserInfoHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserInfoReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.GetUserInfo(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserNameLoginHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserNameLoginReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.UserNameLogin(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserEmailLoginHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserEmailLoginReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.EmailLogin(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserEnableTotpHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserEnableTotpReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.EnableTotp(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserUpdateHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserUpdateReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Update(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserAvatarUpdateHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		avatar, err := c.FormFile("avatar")
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.UpdateAvatar(ctx, avatar)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserGetFriendListHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserGetFriendReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.GetFriendList(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserGetFollowerListHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserGetFollowerReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.GetFollowerList(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserFollowHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserFollowReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Follow(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserUnFollowHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.UserFollowReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.UnFollow(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func UserDeleteHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		l := service.GetUserSrv()
		resp, err := l.Delete(ctx)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
