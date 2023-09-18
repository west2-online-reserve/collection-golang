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

func GetFromContext(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if ok {
		return user, nil
	}
	return nil, errors.New("failed get user_info")
}

func NewContext(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(UserKey).(*UserInfo)
	return user, ok
}
