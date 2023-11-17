package initialize

import (
	"server-redis/cfg"
	"server-redis/midware"
	"server-redis/myredis"
	"server-redis/mysql"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// 数据库初始化
func DatabaseInit() {
	if err:=mysql.ConnectDataBase();err!=nil{
		panic(err)
	}
	if err := mysql.MySQLAccountInit(); err != nil {
		panic(err)
	}
	if err := mysql.MySQLTodolistInit(); err != nil {
		panic(err)
	}
	if err := myredis.RedisInit(); err != nil {
		panic(err)
	}
	if err:=myredis.RedisAccountSync();err!=nil{
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

	s.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("http://"+cfg.ServerHost+"/swagger/doc.json")))

	authorGroup := s.Group("/author", authorizeHandler())

	authorGroup.GET("/ping", authorPingHandler())

	authorGroup.POST("/todolist/add", authorTodolistAddHandler())

	authorGroup.POST("/todolist/search", authorTodolistSearchHandler())

	authorGroup.DELETE("/todolist/delete", authorTodolistDeleteHandler())

	authorGroup.PUT("/todolist/modify", authorTodolistModifyHandler())
}
