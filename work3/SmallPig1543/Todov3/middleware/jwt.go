package middleware

import (
	"Todov3/types"
	"Todov3/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200 //状态码
		token := c.GetHeader("Authorization")
		if token == "" {
			code = http.StatusNotFound
			c.JSON(http.StatusForbidden, gin.H{
				"status":  code,
				"message": "缺少Token",
			})
			c.Abort()
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			code = http.StatusForbidden //无权限
		} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
			code = http.StatusUnauthorized
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status":  code,
				"message": "token有误",
			})
			c.Abort()
			return
		}
		//创建新ctx.request
		c.Request = c.Request.WithContext(types.NewContext(c.Request.Context(), &types.UserInfo{ID: claims.ID}))
		c.Next()
	}
}
