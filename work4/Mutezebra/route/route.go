package route

import (
	"context"
	"fmt"
	"four/api"
	"four/config"
	"four/consts"
	"four/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/http1/resp"
	"net/http"
	"strings"
	"time"
)

func NewRouter() *server.Hertz {
	// 初始化一个实例
	r := server.Default(server.WithHostPorts(config.Config.System.HttpPort), server.WithMaxRequestBodySize(100*consts.MB))
	r.Any("ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "pong")
	})

	// 加载静态资源
	r.StaticFS("/static", &app.FS{Root: "./static"})

	v1 := r.Group("api/")
	{
		// 用户模块
		v1.POST("user/register", api.UserRegisterHandle())
		v1.POST("user/username-login", api.UserNameLoginHandle())
		v1.POST("user/email-login", api.UserEmailLoginHandle())

		// 视频模块
		v1.GET("video/watch", api.VideoWatchContentHandle()) // 看视频
		v1.GET("video/show", api.VideoShowHandle())          // 展示的是视频的基本信息

		// 基础搜索
		v1.POST("search", api.Search())
		v1.POST("search/filter", api.FilterHandle())
		v1.GET("flush/chunk", func(c context.Context, ctx *app.RequestContext) {
			// Hijack the writer of Response
			ctx.Response.HijackWriter(resp.NewChunkedBodyWriter(&ctx.Response, ctx.GetWriter()))

			for i := 0; i < 100; i++ {
				ctx.Write([]byte(fmt.Sprintf("chunk %d: %s\n", i, strings.Repeat("hi~", i)))) // nolint: errcheck
				ctx.Flush()                                                                   // nolint: errcheck
				time.Sleep(1 * time.Second)
			}
		})
	}

	auth := v1.Group("/")
	auth.Use(middleware.JWT()) // 开启jwt鉴权
	{
		// 用户模块
		auth.GET("user/info", api.UserInfoHandle())
		auth.POST("user/enable-totp", api.UserEnableTotpHandle())
		auth.POST("user/update", api.UserUpdateHandle()) //
		auth.POST("user/avatar-update", api.UserAvatarUpdateHandle())
		auth.POST("user/follow", api.UserFollowHandle())
		auth.POST("user/unfollow", api.UserUnFollowHandle())
		auth.GET("user/list/friend", api.UserGetFriendListHandle())     // 双向关注列表
		auth.GET("user/list/follower", api.UserGetFollowerListHandle()) // 关注你的人
		auth.GET("user/delete", api.UserDeleteHandle())

		// 视频模块
		auth.POST("video/comment", api.VideoCommentHandle())
		auth.POST("video/comment/reply", api.VideoCommentReplyHandle())
		auth.POST("video/upload", api.VideoUploadHandle())
		auth.POST("video/delete", api.VideoDeleteHandle())

		// 搜索模块
		auth.GET("auth/historyitem", api.HistorySearchItemsHandle())
		auth.POST("auth/search", api.AuthSearchHandle())
		auth.POST("auth/filter", api.AuthFilterHandle())
	}
	return r
}
