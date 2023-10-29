package main

import (
	"Todov3/DataBase"
	"Todov3/Redis"
	"Todov3/routes"
)

// 接口文档地址:https://apifox.com/apidoc/shared-9f8fb43d-540e-4a1f-85c8-eda709fc6dad
func main() {
	DataBase.LinkMySQL()
	Redis.LinkRedis()
	r := routes.NewRoutes()
	_ = r.Run(":9090")
}
