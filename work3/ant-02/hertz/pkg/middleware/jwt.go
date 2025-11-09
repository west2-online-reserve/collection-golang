package middleware

import (
	"context"
	"hertz/config"
	"hertz/database"
	"hertz/pkg/model"
	"hertz/pkg/repository"
	"hertz/pkg/service"
	"hertz/util"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
)

var (
	once sync.Once
	auth *jwt.HertzJWTMiddleware
)

func NewJWTMiddleware(us service.UserService) *jwt.HertzJWTMiddleware {
	var err error
	identityKey := "uid"
	once.Do(func() {
		cfg := config.GetConfig()
		auth, err = jwt.New(&jwt.HertzJWTMiddleware{
			Key:         []byte(cfg.Jwt.SecretKey),
			Timeout:     time.Hour * cfg.Jwt.ExpireTime,
			MaxRefresh:  time.Hour * cfg.Jwt.MaxRefresh,
			IdentityKey: identityKey,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*model.User); ok {
					return jwt.MapClaims{
						identityKey: v.ID,
					}
				}
				return jwt.MapClaims{}
			},
			IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
				claims := jwt.ExtractClaims(ctx, c)
				if userID, ok := claims[identityKey].(float64); ok {
					return uint64(userID) // 转换为你的目标类型
				}
				return nil
			},
			Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
				var loginVals model.Login
				if err := c.BindAndValidate(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				ussername := loginVals.Username
				password := util.MD5Sum(loginVals.Password)

				us := service.NewUserService(repository.NewUserRepository(database.GetMysql()))
				id := us.Login(ussername, password)
				if id == 0 {
					return nil, jwt.ErrFailedAuthentication
				}

				c.Set(identityKey, id)

				return &model.User{
					ID:       id,
					Username: ussername,
					Password: password,
				}, nil
			},
			Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
				if userID, ok := data.(uint64); ok {
					r := database.GetRedis()
					t1 := jwt.GetToken(ctx, c)
					if t2, err := r.Get(strconv.FormatUint(userID, 10)); err == nil && t1 == t2 {
						return true
					}
				}
				return false
			},
			LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
				// 从上下文中提取用户ID
				uid, ok := c.Get(identityKey)
				if !ok {
					c.JSON(consts.StatusInternalServerError, map[string]interface{}{
						"status": 500,
						"msg":    "登入失败",
					})
				}
				r := database.GetRedis()
				if err := r.Set(strconv.FormatUint(uid.(uint64), 10), token, time.Until(expire)); err != nil {
					c.JSON(consts.StatusInternalServerError, map[string]interface{}{
						"status": 500,
						"msg":    "登入失败",
					})
				}

				// 调用默认的登录响应
				c.JSON(code, map[string]interface{}{
					"status": code,
					"msg":    "登入成功",
					"data": map[string]interface{}{
						"token":  token,
						"expire": expire.Format(time.RFC3339),
					},
				})
			},
		})
	})
	if err != nil {
		log.Fatalf("jwt error: %v", err)
		return nil
	}
	return auth
}
