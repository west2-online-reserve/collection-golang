package api

import (
	"Todov3/service"
	"Todov3/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CreateTaskRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.Create(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func ShowAllTasksHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ShowAllTaskRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.ShowAllTasks(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func ShowAllTasksWithConditionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ShowAllTaskRequestWithCondition
		if err := c.ShouldBind(&req); err == nil && (req.Condition == 0 || req.Condition == 1) {
			var taskServ service.TaskService
			res, err := taskServ.ShowAllTasksWithCondition(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func SearchTasksHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.SearchTasksRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.SearchTasks(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func UpdateTaskHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateTaskRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.UpdateTask(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

// 输入为1，意思是将所有已完成事项改为未完成,反之亦然
func UpdateAllTasksHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UpdateALlTasksRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.UpdateAllTasks(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func DeleteTaskHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.DeleteTaskRequest
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.DeleteTask(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}

func DeleteTasksWithConditionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.DeleteTasksRequestWithCondition
		if err := c.ShouldBind(&req); err == nil {
			var taskServ service.TaskService
			res, err := taskServ.DeleteTasks(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "参数有误",
			})
		}
	}
}
