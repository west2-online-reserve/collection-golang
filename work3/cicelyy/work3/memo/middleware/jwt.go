package middleware

import (
	"memo/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

//返回一个 gin.HandlerFunc 类型的函数
func JWT() gin.HandlerFunc{ 
	return func(c *gin.Context) { 
		code := 200
		//从请求头中获取Authorization字段的值，JWT令牌
		token := c.GetHeader("Authorization") 
		if token == ""{ 
			//设置状态码为 404，表示令牌缺失
			code = 404  
		}else{ 
			//解析令牌
			claim, err := utils.ParseToken(token) 
			//错误处理 
			if err != nil {
				code = 403 //令牌无权限
			}else if time.Now().Unix() > claim.ExpiresAt { 
				code = 401 //令牌无效
			} 
		}
		if code != 200 { 
			c.JSON(200, gin.H{ 
				"status": code,
				"msg": "Token解析错误",
			})
			//中断请求处理流程
			c.Abort()
			return
		}
		//若验证通过，继续处理请求
		c.Next()
	}
}