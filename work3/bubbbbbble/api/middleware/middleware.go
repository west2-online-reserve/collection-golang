package middleware

import (
	"bubbbbbble/config"
	"errors"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
type User struct {
	Name     string `json:"name" `
	Password string `json:"password"`
}
type MyClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}
func GenToken(name string) string {
	key :=[]byte(config.Vp.GetString("jwtkey"))
	claims := MyClaims{
		name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "bubble",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(key)
	return tokenString
}
func ParseToken(tokenString string) (*MyClaims, error) {
	key := []byte(config.Vp.GetString("jwtkey"))
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid Token")
}
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{
				"status": 401,
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(401, gin.H{
				"status": 401,
				"msg":  "请求头中auth格式有误",
			})
			ctx.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			ctx.JSON(401, gin.H{
				"status": 401,
				"msg":  "无效的Token",
			})
			ctx.Abort()
			return
		}
	
		ctx.Set("username", mc.Name)
		ctx.Next()
	}

}

