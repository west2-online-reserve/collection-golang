package main

import (
	"four/config"
	"four/pkg/log"
	"four/pkg/myutils"
	"four/repository/cache"
	"four/repository/db/dao"
	"four/repository/es"
	"four/repository/es/index"
	"four/route"
)

func main() {
	initAll()
	r := route.NewRouter()
	r.Spin()
}

func initAll() {
	config.InitConfig()
	config.DirInit()
	log.InitLog()
	dao.InitMysql()
	cache.InitRedis()
	es.InitES()
	index.InitIndex()
	myutils.OssInit()
}
