package types

import (
	"context"
	"errors"
)

type UserRequest struct {
	UserName string `json:"username" form:"username" binding:"required,max=100"`
	PassWord string `json:"password" form:"password" binding:"required,min=3,max=100"`
}

var userKey int

type UserInfo struct {
	ID uint `json:"id"`
}

// NewContext 新创建一个context，存有userInfo
func NewContext(c context.Context, user *UserInfo) context.Context {
	return context.WithValue(c, userKey, user)
}

// 从context中获取userInfo
func GetUserInfo(c context.Context) (*UserInfo, error) {
	user, ok := c.Value(userKey).(*UserInfo)
	if !ok {
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil
}
