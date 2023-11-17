package main

import (
	"server-redis/initialize"
	"server-redis/cfg"
	_ "server-redis/docs"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// @title test
// @version 1.1-redis
// description test

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 127.0.0.1:8080
// @schemes http
// @BasePath /
func main() {
	initialize.DatabaseInit()
	initialize.MidWareInit()
	s := server.Default(server.WithHostPorts(cfg.ServerHost))
	initialize.ServerInit(s)
	s.Spin()
}

