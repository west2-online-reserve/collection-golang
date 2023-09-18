package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"three/consts"
	"three/pkg/ctl"
	"three/pkg/e"
	"three/pkg/utils"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")
		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  errors.New("token is empty").Error(),
			})
			c.Abort()
			return
		}

		newAToken, newRToken, err := utils.CheckToken(accessToken, refreshToken)
		fmt.Println(accessToken)
		fmt.Println(refreshToken)
		if err != nil {
			code = e.CheckTokenFailed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		claims, err, _ := utils.ParseToken(newAToken)
		if err != nil {
			code = e.ParseTokenFailed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusOK,
				"data":   nil,
				"msg":    e.GetMsg(code),
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		setHeader(c, newAToken, newRToken)
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{ID: claims.ID, UserName: claims.UserName}))
		c.Next()
	}
}

func setHeader(c *gin.Context, aToken, rToken string) {
	c.Header(consts.AccessToken, aToken)
	c.Header(consts.RefreshToken, rToken)
}
