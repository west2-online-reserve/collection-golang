package initialize

import (
	"context"
	"encoding/json"
	"log"
	"server/account"
	"server/datastruct"
	"server/midware"
	"server/mysql"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func registerHandler() app.HandlerFunc {
	return account.Register
}

func loginHandler() app.HandlerFunc {
	return midware.JWTMidWare.LoginHandler
}

func testHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, utils.H{
			"code": consts.StatusOK,
			"msg":  "visited successfully",
		})
	}
}

func authorizeHandler() app.HandlerFunc {
	return midware.JWTMidWare.MiddlewareFunc()
}

func authorPingHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, _ := c.Get(midware.IdentityKey)
		c.JSON(200, utils.H{
			"code":     consts.StatusOK,
			"msg":      "token passed.",
			"username": user.(*datastruct.User).UserName,
		})
	}
}

func authorTodolistAddHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(200, utils.H{
				"code": consts.StatusBadRequest,
				"msg":  err.Error(),
			})
			return
		}
		if isExist {
			var what2do datastruct.TodolistBindJSONReceive
			json.Unmarshal(body, &what2do)
			what2do.Owner = user.(*datastruct.User).UserName
			what2do.Status = false

			if id, err := account.InsertUserTodoList(user.(*datastruct.User).UserName, what2do); err != nil {
				log.Print("[Error] User[", user.(*datastruct.User).UserName, "] operation:", err)
				c.JSON(200, utils.H{
					"code": consts.StatusBadRequest,
					"msg":  err.Error(),
				})
				return
			} else {
				c.JSON(200, utils.H{
					"code":     consts.StatusOK,
					"username": user.(*datastruct.User).UserName,
					"acceptMsg": datastruct.TodolistBindJSONSend{
						Id:       id,
						Title:    what2do.Title,
						Owner:    what2do.Owner,
						Text:     what2do.Text,
						Deadline: what2do.Deadline,
					},
				})
			}
		} else {
			c.JSON(200, utils.H{
				"code": consts.StatusBadRequest,
				"msg":  "identity doesn't exist.",
			})
		}
	}
}

func authorTodolistSearchHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(200, utils.H{
				"code": consts.StatusBadRequest,
				"msg":  err.Error(),
			})
			return
		}
		if isExist {
			var searchCondition datastruct.TodolistBindMysqlSearch
			json.Unmarshal(body, &searchCondition)
			searchCondition.Username = user.(*datastruct.User).UserName
			list, err := mysql.MySQLTodoLIstSearch(searchCondition)

			if err != nil {
				c.JSON(200, utils.H{
					"code": consts.StatusBadRequest,
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, utils.H{
				"code":     consts.StatusOK,
				"username": user.(*datastruct.User).UserName,
				"result":   list,
			})

		} else {
			c.JSON(200, utils.H{
				"code": consts.StatusBadRequest,
				"msg":  "identity doesn't exist.",
			})
		}

	}
}
