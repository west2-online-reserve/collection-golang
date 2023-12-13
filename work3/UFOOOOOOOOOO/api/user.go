package api

import (
	"github.com/gin-gonic/gin"
	"todolist/service"
)

func UserRegister(c *gin.Context){
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister);err == nil{
		res := userRegister.Register()
		c.JSON(200, res)
	}else{
		c.JSON(400, ErrorResponse(err))
	}
}

func Userlogin(c *gin.Context){
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin);err == nil{
		res := userLogin.Login()
		c.JSON(200, res)
	}else{
		c.JSON(400, ErrorResponse(err))
	}
}