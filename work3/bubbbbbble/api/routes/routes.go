package routes

import (
	"github.com/gin-gonic/gin"
	"bubbbbbble/api/controllers"
	"bubbbbbble/api/middleware"
)

func SetRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.SignUp)
	todo := r.Group("/todo", middleware.JWTAuthMiddleware())
	{
		todo.GET("/all", controllers.GetAll)
		todo.GET("/:id", controllers.GetByKey)
		todo.GET("/done", controllers.GetAllDone)
		todo.GET("/undo", controllers.GetAllUndo)
		todo.GET("/all/:key", controllers.GetByKey)
		todo.POST("/alltodo", controllers.Create)
		todo.PUT("/done/:id", controllers.UpdateSingleDone)
		todo.PUT("/undo/:id", controllers.UpdateSingleUndo)
		todo.PUT("/all/done", controllers.UpdateAllDone)
		todo.PUT("/all/undo", controllers.UpdateAllUndo)
		todo.DELETE("/:id", controllers.DelSingle)
		todo.DELETE("/all/done", controllers.DelAllDone)
		todo.DELETE("/all/undo", controllers.DelAllUndo)
		todo.DELETE("/all", controllers.DelAll)
	}
}