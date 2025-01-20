//路由
package routers

import(
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"memo/api"
	"memo/middleware"
)

//初始化并设置 Gin Web 框架的路由
func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/v1"){
		//用户操作
	    v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.Userlogin)
		//身份认证
		authed := v1.Group("/")
		authed.Use(middleware.JWT()) {
			authed.POST("task", api.CreateTask)//增
			authed.GET("task/:id", api.ShowTask)//删
			authed.GET("tasks/", api.ListTask)//查
			authed.PUT("task/:id", api.UpdateTask)
			authed.POST("search", api.SearchTask)//模糊查找
			authed.DELETE("task/:id", api.DeleteTask)
		}
	}
	return r
}