package main

import (
	"bubbbbbble/api/routes"
	"bubbbbbble/config"
	"bubbbbbble/dao"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = dao.InitMysql()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := gin.Default()
	routes.SetRoutes(r)
	serverPort := config.Vp.GetString("serverport")
	r.Run(serverPort)
}
