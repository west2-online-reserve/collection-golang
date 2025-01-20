//使用 Gin Web 框架来处理 HTTP 请求，使用 service 包中的服务来执行业务逻辑，并使用 logrus 库来记录错误日志
//首先解析请求数据，然后执行相应的业务逻辑，最后返回结果或错误
package api

import (
	"memo/service" 
	logging "github.com/sirupsen/logrus" 
	"github.com/gin-gonic/gin" 
	"memo/pkg/utils" 
)

//处理创建任务的请求
func CreateTask (c *gin.Context){
	var createTask service.CreateTaskService
	//使用 utils.ParseToken 从请求头中解析令牌，获取用户 ID
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	//使用 c.ShouldBind 将请求绑定到 service.CreateTaskService 
	if err := c.ShouldBind(&createTask); err == nil {
		res:= createTask.Create(claim.Id) 
		c.JSON(200, res)
	}else{ 
		logging.Error(err) //如果有错返回错误并打印日志 
		c.JSON(400,err)
	}
}

//展示单个任务
func ShowTask(c *gin.Context){
	var showTask service.ShowTaskService
	// claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	//使用 c.ShouldBind 将请求绑定到 service.ShowTaskService
	if err := c.ShouldBind(&showTask); err == nil {
		res:= showTask.Show(c.Param("id"))
		c.JSON(200, res)
	}else{
		logging.Error(err)
		c.JSON(400, err)
	}
}

//处理列出任务的请求
func ListTask(c *gin.Context){
	var listTask service.ListTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&listTask); err == nil {
		res:= listTask.List(claim.Id)
		c.JSON(200, res)
	}else{
		logging.Error(err)
		c.JSON(400,err)
	}
}

//处理更新任务的请求
func UpdateTask(c *gin.Context){
	var updateTask service.UpdateTaskService
	if err := c.ShouldBind(&updateTask); err == nil {
		res:= updateTask.Update(c.Param("id"))
		c.JSON(200, res)
	}else{
		logging.Error(err)
		c.JSON(400, err)
	}
}

//处理搜索任务的请求
func SearchTask(c *gin.Context){
	var searchTask service.SearchTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&searchTask); err == nil {
		res:= searchTask.Search(claim.Id)
		c.JSON(200, res)
	}else{
		logging.Error(err)
		c.JSON(400, err)
	}
}

//处理删除任务的请求
func DeleteTask(c *gin.Context){
	var deleteTask service.DeleteTaskService
	// claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&deleteTask); err == nil {
		res:= deleteTask.Delete(c.Param("id"))
		c.JSON(200, res)
	}else{
		logging.Error(err)
		c.JSON(400, err)
	}
}