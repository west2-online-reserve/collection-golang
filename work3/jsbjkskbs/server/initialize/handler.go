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

// @Summary 注册api
// @Param username query string true "用户名" minlen(1) maxlen(63)
// @Param password query string true "密码"   minlen(6) maxlen(15)
// @Produce application/json
// @Router /register [post]
func registerHandler() app.HandlerFunc {
	return account.Register
}

// @Summary 登录api
// @Desciption 获取token
// @Param username query string true "用户名" minlen(1) maxlen(63)
// @Param password query string true "密码"   minlen(6) maxlen(15)
// @Produce application/json
// @Router /login [post]
func loginHandler() app.HandlerFunc {
	return midware.JWTMidWare.LoginHandler
}

// @Summary 测试api
// @Produce application/json
// @Router /test [get]

func testHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, datastruct.ShortResponse{
			Status:  consts.StatusOK,
			Message: "ok",
			Error:   "",
		})
	}
}

// @Summary token_pass组
// @Description 外部不可用
func authorizeHandler() app.HandlerFunc {
	return midware.JWTMidWare.MiddlewareFunc()
}

// @Summary token测试api
// @Description token前面要添加Bearer
// @Tags author
// @Security ApiKeyAuth
// @Produce application/json
// @Router /author/ping [get]
func authorPingHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, _ := c.Get(midware.IdentityKey)
		c.JSON(consts.StatusOK, datastruct.ShortResponse{
			Status:  consts.StatusOK,
			Message: fmt.Sprintf("token passed. Username: %s", user.(*datastruct.User).UserName),
			Error:   "",
		})
	}
}

// @Summary 添加备忘录api
// @Description token前面要添加Bearer
// @Tags author
// @Security ApiKeyAuth
// @Param data body datastruct.TodolistBindJSONReceive true "标题,内容,截止日期[yyyy-mm-dd hh:mm:ss]"
// @Accept application/json
// @Produce application/json
// @Router /author/todolist/add [post]
func authorTodolistAddHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
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
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(consts.StatusOK, datastruct.SendingJSONFormat{
					Status: consts.StatusOK,
					Data: datastruct.SendingJSONData{
						Items: []datastruct.TodolistBindJSONSend{
							{
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
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
		}
	}
}

// @Summary 查找备忘录api
// @Description token前面要添加Bearer;method(允许叠加,使用或运算):1[isdone],2[keyword],4[all]
// @Tags author
// @Security ApiKeyAuth
// @Param data body datastruct.TodolistBindMysqlSearch true "是否完成,关键字,查找方法"
// @Accept application/json
// @Produce application/json
// @Router /author/todolist/search [post]
func authorTodolistSearchHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()

		if err != nil {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
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
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(consts.StatusOK, datastruct.SendingJSONFormat{
				Status: consts.StatusOK,
				Data: datastruct.SendingJSONData{
					Items:      list,
					TotalItems: len(list),
				},
				Message: "ok",
				Error:   "",
			})

		} else {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
		}

	}
}

// @Summary 删除备忘录api
// @Description token前面要添加Bearer;method(允许叠加,使用或运算):1[isdone],2[idlist],4[all]
// @Tags author
// @Security ApiKeyAuth
// @Param data body datastruct.TodolistBindMysqlDelete true "是否完成,id数组,查找方法"
// @Accept application/json
// @Produce application/json
// @Router /author/todolist/delete [delete]
func authorTodolistDeleteHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()
		if err != nil {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
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
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(consts.StatusOK, datastruct.ShortResponse{
				Status:  consts.StatusOK,
				Message: "ok",
				Error:   "",
			})
		} else {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
	}
}

// @Summary 更新备忘录api
// @Description token前面要添加Bearer;method(允许叠加,使用或运算):1[isdone],2[idlist],4[all]
// @Tags author
// @Security ApiKeyAuth
// @Param data body datastruct.TodolistBindMysqlModify true "是否完成,id数组,查找方法"
// @Accept application/json
// @Produce application/json
// @Router /author/todolist/modify [put]
func authorTodolistModifyHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()
		if err != nil {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
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
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(consts.StatusOK, datastruct.ShortResponse{
				Status:  consts.StatusOK,
				Message: "ok",
				Error:   "",
			})
		}
	}
}
