package main

import (
	"todolist/conf"
	"todolist/routers"
)

func main() {
	conf.Init()
	r := routers.NewRouter()
	_=r.Run(conf.HttpPort)
}