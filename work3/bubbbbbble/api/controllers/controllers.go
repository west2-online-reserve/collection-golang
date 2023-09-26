package controllers

import (
	"bubbbbbble/api/middleware"
	"bubbbbbble/model"
	"bubbbbbble/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(ctx *gin.Context) {
	name := ctx.GetString("username")
	todolist, err := service.GetAll(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"data": gin.H{
				"items": "null",
				"total": 0,
			},
			"msg":   "Not found",
			"error": "",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "200",
			"data": gin.H{
				"items": todolist,
				"total": len(todolist),
			},
			"msg":   "ok",
			"error": "",
		})
	}
}
func GetAllDone(ctx *gin.Context) {
	name := ctx.GetString("username")
	todolist, err := service.GetAllDone(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"data": gin.H{
				"items": todolist,
				"total": len(todolist),
			},
			"msg":   "Not found",
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"items": todolist,
			"total": len(todolist),
		},
		"msg":   "ok",
		"error": "",
	})
}
func GetAllUndo(ctx *gin.Context) {
	name := ctx.GetString("username")
	todolist, err := service.GetAllUndo(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"data": gin.H{
				"items": todolist,
				"total": len(todolist),
			},
			"msg":   "Not found",
			"error": err,
		})
		return
	}
	for _, todo := range todolist {
		todo.Viewed++
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"items": todolist,
			"total": len(todolist),
		},
		"msg":   "ok",
		"error": "",
	})
}
func GetSingle(ctx *gin.Context) {
	var todo model.Todo
	ctx.Param("id")
	id, _ := strconv.Atoi(ctx.Param("id"))
	name := ctx.GetString("username")
	todo, err := service.GetSingle(name, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"data": gin.H{
				"items": todo,
				"total": 0,
			},
			"msg":   "Not found",
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"items": todo,
			"total": 1,
		},
		"msg":   "ok",
		"error": "",
	})

}
func GetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	name := ctx.GetString("username")
	todolist, err := service.GetByKey(name, key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"data": gin.H{
				"items": todolist,
				"total": len(todolist),
			},
			"msg":   "Not found",
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"items": todolist,
			"total": len(todolist),
		},
		"msg":   "ok",
		"error": "",
	})
}
func Create(ctx *gin.Context) {
	var todo model.Todo
	ctx.ShouldBind(&todo)
	name := ctx.GetString("username")
	if name != todo.Name {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":  "401",
			"msg":   "failure",
			"error": "操作权限不足",
		})
		return
	}
	service.Create(todo)
	ctx.JSON(http.StatusCreated, gin.H{
		"code":  "201",
		"msg":   "success",
		"error": "",
	})
}
func UpdateAllDone(ctx *gin.Context) {
	name := ctx.GetString("username")
	service.UpdateAllDone(name)
	ctx.JSON(http.StatusCreated, gin.H{
		"code":  "201",
		"msg":   "success",
		"error": "",
	})

}
func UpdateAllUndo(ctx *gin.Context) {
	name := ctx.GetString("username")
	service.UpdateAllUndo(name)
	ctx.JSON(http.StatusCreated, gin.H{
		"code":  "201",
		"msg":   "success",
		"error": "",
	})
}
func UpdateSingleDone(ctx *gin.Context) {
	name := ctx.GetString("username")
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := service.UpdateSingleDone(name, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":  "404",
			"msg":   "Not found",
			"error": "id未找到",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code":  "201",
		"msg":   "success",
		"error": "",
	})
}
func UpdateSingleUndo(ctx *gin.Context) {
	name := ctx.GetString("username")
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := service.UpdateSingleUndo(name, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":  "404",
			"msg":   "Not found",
			"error": "id不存在",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code":  "201",
		"msg":   "success",
		"error": "",
	})
}

func DelAllDone(ctx *gin.Context) {
	name := ctx.GetString("username")
	service.DelAllDone(name)
	ctx.JSON(http.StatusNoContent, gin.H{
		"code": 204,
		"msg":  "success",
	})
}
func DelAllUndo(ctx *gin.Context) {
	name := ctx.GetString("username")
	service.DelAllUndo(name)
	ctx.JSON(http.StatusNoContent, gin.H{
		"code": 204,
		"msg":  "success",
	})
}
func DelSingle(ctx *gin.Context) {
	name := ctx.GetString("username")
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := service.DelSingle(name, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"msg":  "Not found",
		})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{
		"code": 204,
		"msg":  "success",
	})
}
func DelAll(ctx *gin.Context) {
	name := ctx.GetString("username")
	service.DelAll(name)
	ctx.JSON(http.StatusNoContent, gin.H{
		"code": 204,
		"msg":  "success",
	})
}
func SignUp(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBind(&user)
	err := service.SignUp(&user)
	if err != nil {
		service.CreateUser(user)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
		})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "registered",
		})

	}
}
func Login(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBind(&user)
	err := service.Login(&user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 400,
			"msg":  "用户不存在或密码不匹配",
		})
	} else {
		tokenString := middleware.GenToken(user.Name)
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"code": 200,
			"data": gin.H{
				"token": tokenString,
			},
		})
	}
}
