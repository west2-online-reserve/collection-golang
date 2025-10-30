package service

import (
	"errors"
	"memogo/biz/dal/model"
	"memogo/biz/dal/repository"
	"memogo/pkg/hash"
	jwtPkg "memogo/pkg/jwt"
)

var (
	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUsernameRequired 用户名不能为空
	ErrUsernameRequired = errors.New("username is required")
	// ErrPasswordRequired 密码不能为空
	ErrPasswordRequired = errors.New("password is required")
	// ErrPasswordTooShort 密码太短
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *AuthService) Register(username, password string) (accessToken, refreshToken string, err error) {
	// 参数验证
	if username == "" {
		return "", "", ErrUsernameRequired
	}
	if password == "" {
		return "", "", ErrPasswordRequired
	}
	if len(password) < 6 {
		return "", "", ErrPasswordTooShort
	}

	// 检查用户名是否已存在
	_, err = s.userRepo.GetByUsername(username)
	if err == nil {
		return "", "", repository.ErrUserAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return "", "", err
	}

	// 对密码进行哈希
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return "", "", err
	}

	// 创建用户
	user := &model.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", "", err
	}

	// 生成 JWT 令牌对
	accessToken, refreshToken, err = jwtPkg.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (accessToken, refreshToken string, err error) {
	// 参数验证
	if username == "" {
		return "", "", ErrUsernameRequired
	}
	if password == "" {
		return "", "", ErrPasswordRequired
	}

	// 查找用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	// 验证密码
	if err := hash.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// 生成 JWT 令牌对
	accessToken, refreshToken, err = jwtPkg.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// 解析刷新令牌
	claims, err := jwtPkg.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 验证用户是否存在
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", err
	}

	// 生成新的令牌对
	newAccessToken, newRefreshToken, err = jwtPkg.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
