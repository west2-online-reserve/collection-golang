package routes

import (
	"Todov3/api"
	"Todov3/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRoutes() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret")) // 存储
	r.Use(sessions.Sessions("mysession", store))
	v3 := r.Group("api/v3")
	{
		v3.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})
		v3.POST("user/register", api.RegisterHandler())
		v3.POST("user/login", api.LoginHandler())
		task := v3.Group("task")
		task.Use(middleware.JWT())
		{
			//增
			task.POST("create", api.CreateHandler())
			//查
			task.POST("show_all_tasks", api.ShowAllTasksHandler())
			task.POST("show_all_tasks_with_condition", api.ShowAllTasksWithConditionHandler())
			task.POST("search", api.SearchTasksHandler())
			//改
			task.PUT("update_task", api.UpdateTaskHandler())
			task.PUT("update_tasks", api.UpdateAllTasksHandler())
			//删
			task.DELETE("delete", api.DeleteTaskHandler())
			task.DELETE("delete_with_condition", api.DeleteTasksWithConditionHandler())
		}

	}
	return r
}
