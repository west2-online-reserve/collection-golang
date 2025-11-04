package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 无效的 token
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken 过期的 token
	ErrExpiredToken = errors.New("token expired")
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	OrigIat  int64  `json:"orig_iat"`
	jwt.RegisteredClaims
}

// GetJWTSecret 获取 JWT 密钥
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 警告：生产环境必须设置环境变量
		secret = "memogo-default-secret-change-in-production"
	}
	return []byte(secret)
}

// GenerateToken 生成 JWT token
func GenerateToken(userID uint, username string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		OrigIat:  now.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func GenerateTokenPair(userID uint, username string) (accessToken, refreshToken string, err error) {
	// 访问令牌：15分钟
	accessToken, err = GenerateToken(userID, username, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	// 刷新令牌：7天
	refreshToken, err = GenerateToken(userID, username, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
