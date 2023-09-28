package main

import (
	"server/initialize"
	"server/cfg"
	_ "server/docs"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// @title test
// @version 1.0
// description test

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 127.0.0.1:8080
// @schemes http
// @BasePath /
func main() {
	initialize.MySQLInit()
	initialize.MidWareInit()
	s := server.Default(server.WithHostPorts(cfg.ServerHost))
	initialize.ServerInit(s)
	s.Spin()
}

