package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/account"
	"server/datastruct"
	"server/midware"
	"server/mysql"

	"github.com/cloudwego/hertz/pkg/app"
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
		c.JSON(200, datastruct.ShortResponse{
			Status:  consts.StatusOK,
			Message: "ok",
			Error:   "",
		})
	}
}

func authorizeHandler() app.HandlerFunc {
	return midware.JWTMidWare.MiddlewareFunc()
}

func authorPingHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, _ := c.Get(midware.IdentityKey)
		c.JSON(200, datastruct.ShortResponse{
			Status:  consts.StatusBadRequest,
			Message: fmt.Sprintf("token passed. Username: %s", user.(*datastruct.User).UserName),
			Error:   "",
		})
	}
}

func authorTodolistAddHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
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
				c.JSON(200, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(200, datastruct.SendingJSONFormat{
					Status: consts.StatusOK,
					Data: datastruct.SendingJSONData{
						Items: []datastruct.TodolistBindJSONSend{
							datastruct.TodolistBindJSONSend{
								Id:       id,
								Title:    what2do.Title,
								Owner:    what2do.Owner,
								Text:     what2do.Text,
								Addtime:  "just now",
								Deadline: what2do.Deadline,
								Status:   false,
							},
						},
						TotalItems: 1,
					},
					Message: "ok",
					Error:   "",
				})
			}
		} else {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
		}
	}
}

func authorTodolistSearchHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
		if isExist {
			var searchCondition datastruct.TodolistBindMysqlSearch
			json.Unmarshal(body, &searchCondition)
			searchCondition.Username = user.(*datastruct.User).UserName
			list, err := mysql.MySQLTodoListSearch(searchCondition)

			if err != nil {
				c.JSON(200, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}

			c.JSON(200, datastruct.SendingJSONFormat{
				Status: consts.StatusOK,
				Data: datastruct.SendingJSONData{
					Items:      list,
					TotalItems: len(list),
				},
				Message: "ok",
				Error:   "",
			})

		} else {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
		}

	}
}

func authorTodolistDeleteHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()
		if err != nil {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
		if isExist {
			var deleteCondition datastruct.TodolistBindMysqlDelete
			json.Unmarshal(body, &deleteCondition)
			deleteCondition.Username = user.(*datastruct.User).UserName
			if err := mysql.MySQLTodoListDelete(deleteCondition); err != nil {
				c.JSON(200, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusOK,
				Message: "ok",
				Error:   "",
			})
		} else {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
	}
}

func authorTodolistModifyHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()
		if err != nil {
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
		if isExist {
			var modifyCondition datastruct.TodolistBindMysqlModify
			json.Unmarshal(body, &modifyCondition)
			modifyCondition.Username = user.(*datastruct.User).UserName

			if err := mysql.MySQLTodoListModify(modifyCondition); err != nil {
				c.JSON(200, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(200, datastruct.ShortResponse{
				Status:  consts.StatusOK,
				Message: "ok",
				Error:   "",
			})
		}
	}
}
