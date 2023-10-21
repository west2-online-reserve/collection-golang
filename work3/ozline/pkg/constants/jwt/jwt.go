package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants"
)

type Claims struct {
	UserID int64 `json:"userid"`
	jwt.StandardClaims
}

func CreateToken(userID int64) (string, error) {
	expireTime := time.Now().Add(24 * 7 * time.Hour)
	now := time.Now()

	claims := Claims{
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

func CheckToken(token string) (*Claims, error) {
	response, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
