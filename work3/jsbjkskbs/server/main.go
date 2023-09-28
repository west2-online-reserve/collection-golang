package main

import (
	"server/initialize"
	"server/cfg"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	initialize.MySQLInit()
	initialize.MidWareInit()
	s := server.Default(server.WithHostPorts(cfg.ServerHost))
	initialize.ServerInit(s)
	s.Spin()
}
