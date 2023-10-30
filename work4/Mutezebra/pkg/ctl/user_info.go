package ctl

import (
	"context"
	"errors"
)

type key int64

var UserKey key

type UserInfo struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
}

func NewContext(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetFromContext(ctx context.Context) (*UserInfo, error) {
	userInfo, ok := FromContext(ctx)
	if ok {
		return userInfo, nil
	}
	return nil, errors.New("get user info failed")
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	userInfo, ok := ctx.Value(UserKey).(*UserInfo)
	return userInfo, ok
}
