package initialize

import (
	"server/midware"
	"server/mysql"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// 数据库初始化
func MySQLInit() {
	err := mysql.ConnectDataBase()
	if err != nil {
		panic(err)
	}
	err = mysql.MySQLAccountInit()
	if err != nil {
		panic(err)
	}
	err = mysql.MySQLTodolistInit()
	if err !=nil{
		panic(err)
	}
}

// 中间件初始化
func MidWareInit() {
	midware.JWTInit()
}

// 服务端(router?)初始化
func ServerInit(s *server.Hertz) {

	s.POST("/register", registerHandler())

	s.POST("/login", loginHandler())

	s.GET("/test", testHandler())

	authorGroup := s.Group("/author", authorizeHandler())

	authorGroup.GET("/ping", authorPingHandler())

	authorGroup.POST("/todolist/add", authorTodolistAddHandler())

	authorGroup.POST("/todolist/search", authorTodolistSearchHandler())
}
