package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"three/pkg/utils"
	"three/service"
	"three/types"
)

func UserRegisterHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserRegisterReq
		err := c.ShouldBind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Register(c.Request.Context(), &req)
			if err != nil {
				utils.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		} else {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
	}
}

func UserLoginHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserLoginReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		l := service.GetUserSrv()
		resp, err := l.Login(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}
