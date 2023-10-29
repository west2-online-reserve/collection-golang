package util

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var mySigningKey = []byte("AllYourBase")

type Claims struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GenerateToken(id uint, username string) (string, error) {
	claims := Claims{
		ID:       id,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "Todov3",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(mySigningKey)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
