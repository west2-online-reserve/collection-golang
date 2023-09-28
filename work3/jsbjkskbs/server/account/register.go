package account

import (
	"context"
	"log"
	"server/mysql"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Register(ctx context.Context, c *app.RequestContext) {
	var reigisterStruct struct {
		Username string `form:"username" json:"username" vd:"(len($)>0&&len($)<64); msg:'Illegal Username'"`
		Password string `form:"password" json:"password" vd:"(len($)>5&&len($)<16); msg:'Illegal Password'"`
	}

	if err := c.BindAndValidate(&reigisterStruct); err != nil {
		c.JSON(200, utils.H{
			"code": consts.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	accountList, err := mysql.MySQLAccountSearch(reigisterStruct.Username)

	if err != nil {
		c.JSON(200, utils.H{
			"code": consts.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	if len(accountList) != 0 {
		c.JSON(200, utils.H{
			"code": consts.StatusBadRequest,
			"msg":  "username exists",
		})
		return
	}

	err = mysql.MySQLAccountCreate(reigisterStruct.Username, reigisterStruct.Password)

	if err != nil {
		c.JSON(200, utils.H{
			"code": consts.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, utils.H{
		"code": consts.StatusOK,
		"msg":  "registered successfully",
	})

	log.Printf("[INFO] User [%s] has registered successfully.", reigisterStruct.Username)
}
