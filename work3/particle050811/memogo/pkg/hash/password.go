package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost bcrypt 默认成本因子
	DefaultCost = 12
)

// HashPassword 使用 bcrypt 对密码进行哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
