package api

import (
	"github.com/gin-gonic/gin"
	"memo/service"
)  

//用户注册接口 
func UserRegister(c *gin.Context) {
	var userRegister service.UserService //声明user用户服务对象
	if err := c.ShouldBind(&userRegister);err == nil { //绑定
		res := userRegister.Register() //方法
		c.JSON(200, res) 
	}else { 
		c.JSON(400, err)
	}
} 

//处理用户登录的请求
func Userlogin(c *gin.Context){
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin);err == nil{
		res := userLogin.Login()
		c.JSON(200, res)
	}else{
		c.JSON(400,err)
	} 
} 