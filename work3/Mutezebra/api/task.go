package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"three/pkg/utils"
	"three/service"
	"three/types"
)

func TaskCreateHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskCreateReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.Create(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func TaskUpdateHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskUpdateReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.Update(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func TaskShowHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskShowReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.Show(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func TaskListHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskListReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.List(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func TaskSearchHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskSearchReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		switch req.Text {
		case "":
			resp, err := l.SearchByStatus(c.Request.Context(), &req)
			if err != nil {
				utils.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		default:
			resp, err := l.SearchByText(c.Request.Context(), &req)
			if err != nil {
				utils.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
	}
}

func TaskDeleteHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskDeleteReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.Delete(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func TaskDeleteAllHandel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TaskDeleteReq
		err := c.ShouldBind(&req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		l := service.GetTaskSrv()
		resp, err := l.DeleteAllTask(c.Request.Context(), &req)
		if err != nil {
			utils.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}
