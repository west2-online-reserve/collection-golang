package main

import (
	"three/config"
	"three/pkg/utils"
	"three/repository/cache"
	"three/repository/db/dao"
	"three/routes"
)

func main() {
	initAll()
	r := routes.NewRouter()
	_ = r.Run(config.Config.System.HttpPort)
}

func initAll() {
	utils.InitLog()
	config.InitConfig()
	cache.InitRedis()
	dao.InitMysql()
}
