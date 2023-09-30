package midware

import (
	"context"
	"log"
	"server-redis/datastruct"
	encrypt2 "server-redis/encrypt"
	"server-redis/myredis"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var (
	JWTMidWare  *jwt.HertzJWTMiddleware
	IdentityKey = "indentity"
)

// 密钥
const secretKey = "114514"

func JWTInit() {
	var err error
	JWTMidWare, err = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(secretKey),                                  //签名密钥
		Timeout:       time.Hour,                                          //token有效时间
		MaxRefresh:    time.Hour,                                          //token最大刷新时间
		TokenLookup:   "header: Authorization, query: token, cookie: jwt", //token获取源 [header,query,cookie,param,form]
		TokenHeadName: "Bearer",                                           //header中的token前缀
		IdentityKey:   IdentityKey,                                        //检索身份的key

		//SigningAlgorithm: "HS256",														//加密算法(optional)

		//获取身份信息
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			return &datastruct.User{Username: jwt.ExtractClaims(ctx, c)[IdentityKey].(string)}
		},

		//认证用户登录信息
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" vd:"(len($)>0&&len($)<64); msg:'Illegal Username'"`
				Password string `form:"password" json:"password" vd:"(len($)>5&&len($)<16); msg:'Illegal Password'"`
			}

			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}

			_, err := myredis.RedisCheckAccount(loginStruct.Username, encrypt2.SHA256(loginStruct.Password))

			if err != nil {
				return nil, err
			}

			log.Printf("[INFO] User [%s] has logined successfully.", loginStruct.Username)

			return loginStruct, nil
		},

		//登入相应
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(200, utils.H{
				"status":  code,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "ok",
				"error":   "",
			})
		},

		//登出响应
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			c.JSON(200, datastruct.ShortResponse{
				Status:  code,
				Message: "logout successfully",
				Error:   "",
			})
		},

		//设置 JWT token验证错误的响应信息
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt error:", e.Error())
			return e.Error()
		},

		//JWT token验证失败的响应
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, datastruct.ShortResponse{
				Status:  code,
				Message: "token claiming failed",
				Error:   message,
			})
		},

		//使用JWT token时不可去掉
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if username := data.(struct {
				Username string "form:\"username\" json:\"username\" vd:\"(len($)>0&&len($)<64); msg:'Illegal Username'\""
				Password string "form:\"password\" json:\"password\" vd:\"(len($)>5&&len($)<16); msg:'Illegal Password'\""
			}).Username; len(username) != 0 {
				return jwt.MapClaims{
					IdentityKey: username,
				}
			}
			return jwt.MapClaims{}
		},
	})
	if err != nil {
		panic(err)
	}
}
