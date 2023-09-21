package main

import (
	"four/config"
	"four/pkg/log"
	"four/repository/cache"
	"four/repository/db/dao"
)

func main() {
	initAll()
}

func initAll() {
	config.InitConfig()
	log.InitLog()
	dao.InitMysql()
	cache.InitRedis()
}
