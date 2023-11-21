package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server-redis/account"
	"server-redis/cfg"
	"server-redis/dataprocesser"
	"server-redis/datastruct"
	"server-redis/midware"
	"strconv"

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
			Message: fmt.Sprintf("token passed. Username: %s", user.(*datastruct.User).Username),
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
			what2do.Owner = user.(*datastruct.User).Username
			what2do.Status = false

			if total, id, err := dataprocesser.InsertReceiveTodolist(what2do); err != nil {
				log.Print("[Error] User[", user.(*datastruct.User).Username, "] operation:", err)
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
						Count: 1,
						Total: int(total),
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
// @Param page query int false "页码"
// @Param data body datastruct.TodolistBindRedisCondition true "是否完成,关键字,idlist不填,查找方法"
// @Accept application/json
// @Produce application/json
// @Router /author/todolist/search [post]
func authorTodolistSearchHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		user, isExist := c.Get(midware.IdentityKey)
		body, err := c.Body()
		page, _ := c.GetQuery("page")
		pageNumber, _ := strconv.ParseInt(page, 10, 64)
		if pageNumber == 0 {
			pageNumber = 1
		}
		if err != nil {
			c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
				Status:  consts.StatusBadRequest,
				Message: "",
				Error:   err.Error(),
			})
			return
		}
		if isExist {
			var searchCondition datastruct.TodolistBindRedisCondition
			json.Unmarshal(body, &searchCondition)
			data, isend, err := dataprocesser.SearchUserTodoList(user.(*datastruct.User).Username, searchCondition, int(pageNumber))

			if err != nil {
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			}
			c.JSON(consts.StatusOK, datastruct.SendingJSONFormat{
				Status:       consts.StatusOK,
				Page:         int(pageNumber),
				ItemsPerPage: cfg.ItemsCountInPage,
				Isend:        isend,
				Data:         data,
				Message:      "ok",
				Error:        "",
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
// @Param data body datastruct.TodolistBindRedisCondition true "是否完成,keyword不填,id数组,查找方法"
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
			var deleteCondition datastruct.TodolistBindRedisCondition
			json.Unmarshal(body, &deleteCondition)
			if err := dataprocesser.DeleteUserTodoList(user.(*datastruct.User).Username, deleteCondition); err != nil {
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(consts.StatusOK, datastruct.ShortResponse{
					Status:  consts.StatusOK,
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
			return
		}
	}
}

// @Summary 更新备忘录api
// @Description token前面要添加Bearer;method(允许叠加,使用或运算):2[idlist],4[all]
// @Tags author
// @Security ApiKeyAuth
// @Param data body datastruct.TodolistBindRedisUpdate true "id数组,更新状态,查找方法"
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
			var modifyCondition datastruct.TodolistBindRedisUpdate
			json.Unmarshal(body, &modifyCondition)

			if err := dataprocesser.UpdateTodoList(user.(*datastruct.User).Username, modifyCondition); err != nil {
				c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
					Status:  consts.StatusBadRequest,
					Message: "",
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(consts.StatusOK, datastruct.ShortResponse{
					Status:  consts.StatusOK,
					Message: "ok",
					Error:   "",
				})
			}
		}
	}
}
