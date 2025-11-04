package middleware

import (
	"context"
	"errors"
	"memogo/biz/dal/db"
	"memogo/biz/dal/repository"
	"memogo/pkg/hash"
	"memogo/pkg/jwt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzJWT "github.com/hertz-contrib/jwt"
)

var (
	// JWTMiddleware Hertz JWT 中间件实例
	JWTMiddleware *hertzJWT.HertzJWTMiddleware
	// IdentityKey 用户身份标识的 key
	IdentityKey = "user_id"
)

// JWTClaims JWT 自定义声明
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// InitJWTMiddleware 初始化 JWT 中间件
func InitJWTMiddleware() (*hertzJWT.HertzJWTMiddleware, error) {
	authMiddleware, err := hertzJWT.New(&hertzJWT.HertzJWTMiddleware{
		Realm:       "memogo",
		Key:         jwt.GetJWTSecret(),
		Timeout:     15 * time.Minute,  // Access token 过期时间
		MaxRefresh:  7 * 24 * time.Hour, // Refresh token 过期时间
		IdentityKey: IdentityKey,

		// PayloadFunc 设置 JWT 的 payload
		PayloadFunc: func(data interface{}) hertzJWT.MapClaims {
			if v, ok := data.(*JWTClaims); ok {
				return hertzJWT.MapClaims{
					"user_id":  v.UserID,
					"username": v.Username,
					"orig_iat": time.Now().Unix(), // 原始签发时间，刷新 token 时需要
				}
			}
			return hertzJWT.MapClaims{}
		},

		// IdentityHandler 从 token 中提取用户身份信息
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := hertzJWT.ExtractClaims(ctx, c)
			userID := uint(claims["user_id"].(float64))
			username := claims["username"].(string)

			return &JWTClaims{
				UserID:   userID,
				Username: username,
			}
		},

		// Authenticator 验证用户登录
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginReq struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.BindAndValidate(&loginReq); err != nil {
				return nil, err
			}

			// 调用认证逻辑
			userRepo := repository.NewUserRepository(db.DB)
			user, err := userRepo.GetByUsername(loginReq.Username)
			if err != nil {
				return nil, hertzJWT.ErrFailedAuthentication
			}

			// 验证密码
			if err := hash.VerifyPassword(user.PasswordHash, loginReq.Password); err != nil {
				return nil, hertzJWT.ErrFailedAuthentication
			}

			return &JWTClaims{
				UserID:   user.ID,
				Username: user.Username,
			}, nil
		},

		// LoginResponse 登录成功后的响应
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(200, utils.H{
				"status": 200,
				"msg":    "Login successful",
				"data": utils.H{
					"access_token":  token,
					"refresh_token": token, // 这里暂时返回相同的 token
					"expires_at":    expire.Unix(),
				},
			})
		},

		// RefreshResponse 刷新 token 成功后的响应
		RefreshResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(200, utils.H{
				"status": 200,
				"msg":    "Token refreshed",
				"data": utils.H{
					"access_token":  token,
					"refresh_token": token,
					"expires_at":    expire.Unix(),
				},
			})
		},

		// Unauthorized 未授权的响应
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, utils.H{
				"status": code,
				"msg":    message,
				"data":   nil,
			})
		},

		// TokenLookup 从请求中查找 token 的位置
		TokenLookup: "header: Authorization",

		// TokenHeadName token 的前缀
		TokenHeadName: "Bearer",

		// TimeFunc 时间函数
		TimeFunc: time.Now,
	})

	if err != nil {
		return nil, err
	}

	// 初始化中间件
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		return nil, errors.New("JWT middleware initialization failed: " + errInit.Error())
	}

	return authMiddleware, nil
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *app.RequestContext) (uint, error) {
	claims, exists := c.Get(IdentityKey)
	if !exists {
		return 0, errors.New("user not found in context")
	}

	jwtClaims, ok := claims.(*JWTClaims)
	if !ok {
		return 0, errors.New("invalid claims type")
	}

	return jwtClaims.UserID, nil
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *app.RequestContext) (string, error) {
	claims, exists := c.Get(IdentityKey)
	if !exists {
		return "", errors.New("user not found in context")
	}

	jwtClaims, ok := claims.(*JWTClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	return jwtClaims.Username, nil
}
