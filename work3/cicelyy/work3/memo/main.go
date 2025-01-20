package main

import (
	"memo/conf"
	"memo/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	_=r.Run(conf.HttpPort)  
}