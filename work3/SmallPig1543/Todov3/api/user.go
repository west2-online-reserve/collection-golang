package api

import (
	"Todov3/service"
	"Todov3/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserRequest
		if err := c.ShouldBind(&req); err == nil {
			var userServ service.UserService
			res, err := userServ.Register(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserRequest
		if err := c.ShouldBind(&req); err == nil {
			var userServ service.UserService
			res, err := userServ.Login(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	}
}
