package jwt

import (
	"Demo/pkg/constants"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(userID int64) (string, error) {
	expireTime := time.Now().Add(24 * 7 * time.Hour)
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    constants.JWTValue,
			IssuedAt:  now.Unix(),
		},
	}

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenStruct.SignedString([]byte(constants.JWTValue))
}

func CheckToken(tokenString string) (*Claims, error) {
	response, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTValue), nil
	})
	if err != nil {
		return nil, err
	}
	if resp, ok := response.Claims.(*Claims); ok && response.Valid {
		return resp, nil
	}
	return nil, err
}
