package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var JWTsecret = []byte("ABAB") 

//定义 JWT 中的声明
type Claims struct { 
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	jwt.StandardClaims 
}

//签发token 用于生成 JWT
func GenerateToken (id uint, username, password string) (string, error){
	//时间戳
	notTime := time.Now()
	//设置过期时间 
	expireTime := notTime.Add(3*time.Hour)
	claims := Claims{
		Id : id,
		UserName: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : expireTime.Unix(),
			Issuer : "memo",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTsecret)
	return token, err 
}

//验证token 解析 JWT
func ParseToken(token string) (*Claims, error){
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, 
	func(token *jwt.Token) (interface{}, error){
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}