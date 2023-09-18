package routes

import (
	"github.com/gin-gonic/gin"
	"three/api"
	"three/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.String(200, "success")
		})
		v1.POST("user/register", api.UserRegisterHandel())
		v1.POST("user/login", api.UserLoginHandel())

		auth := v1.Group("/")
		auth.Use(middleware.JWT())
		{
			auth.POST("task/create", api.TaskCreateHandel())
			auth.PUT("task/update", api.TaskUpdateHandel())
			auth.GET("task/show", api.TaskShowHandel())
			auth.GET("task/list", api.TaskListHandel())
			auth.POST("task/search", api.TaskSearchHandel())
			auth.DELETE("task/delete", api.TaskDeleteHandel())
			auth.DELETE("task/deleteAll", api.TaskDeleteAllHandel())
		}
	}
	return r
}
