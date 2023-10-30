package myutils

import (
	"four/consts"
	"four/pkg/log"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserName string
	ID       uint
	jwt.StandardClaims
}

// CheckToken 用于jwt中间件里检查token
func CheckToken(aToken, rToken string) (newAToken, newRToken string, err error) {

	aClaims, err, aValid := ParseToken(aToken)
	if err != nil {
		log.LogrusObj.Println("atoken", err)
		return
	}
	rClaims, err, rValid := ParseToken(rToken)
	if err != nil {
		log.LogrusObj.Println("rtoken", err)
		return
	}
	// 如果aToken和rToken都不过期就只更新aToken
	if rValid && aValid {
		newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.ID)
		newRToken = rToken
		return
	}
	// 如果aToken过期 但是rToken不过期就只更新aToken
	if rValid && !aValid {
		newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.ID)
		newRToken = rToken
		return
	}
	// 全更新
	newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.ID)
	if err != nil {
		return
	}
	newRToken, err = GenerateRefreshToken(rClaims.UserName, rClaims.ID)
	if err != nil {
		return
	}
	return
}

// GenerateToken 登陆时签发Token
func GenerateToken(userName string, id uint) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userName, id)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(userName, id)
	if err != nil {
		return "", "", err
	}
	return
}

// GenerateAccessToken 签发AccessToken
func GenerateAccessToken(userName string, id uint) (accessToken string, err error) {
	timeNow := time.Now()
	accessTokenExpireTime := timeNow.Add(consts.AccessTokenExpireTime).Unix()
	claims := &Claims{
		UserName: userName,
		ID:       id,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Mutezebra",
			Subject:   userName,
			ExpiresAt: accessTokenExpireTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(consts.JwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GenerateRefreshToken 签发AccessToken
func GenerateRefreshToken(userName string, id uint) (refreshToken string, err error) {
	timeNow := time.Now()
	refreshTokenExpireTime := timeNow.Add(consts.RefreshTokenExpireTime).Unix()
	claims := &Claims{
		UserName: userName,
		ID:       id,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Mutezebra",
			Subject:   userName,
			ExpiresAt: refreshTokenExpireTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err = token.SignedString([]byte(consts.JwtSecret))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// ParseToken 解析token并判断其有没有过期
func ParseToken(token string) (*Claims, error, bool) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtSecret), nil
	})

	if err != nil {
		return nil, err, false
	}
	claims, ok := tokenClaims.Claims.(*Claims)
	if ok && tokenClaims.Valid {
		return claims, nil, IsValid(tokenClaims)
	}
	return nil, err, false
}

func IsValid(token *jwt.Token) bool {
	return token.Valid
}
