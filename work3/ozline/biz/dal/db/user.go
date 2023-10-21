package db

import (
	"context"
	"errors"

	"github.com/west2-online-reserve/collection-golang/work3/pkg/constants"
)

func CreateUser(ctx context.Context, username, password string) (*User, error) {
	var userResp *User

	err := DB.WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", username).
		First(&userResp).
		Error

	if err == nil {
		return nil, errors.New("user existed")
	}

	userResp = &User{
		Username: username,
		Password: password,
	}

	err = DB.WithContext(ctx).
		Table(constants.UserTable).
		Create(&userResp).
		Error

	if err != nil {
		return nil, err
	}

	return userResp, nil
}

func LoginCheck(ctx context.Context, username, password string) (*User, error) {
	var userResp *User

	err := DB.WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", username).
		First(&userResp).
		Error

	if err != nil {
		return nil, err
	}

	if userResp.Password == password {
		return userResp, nil
	}

	return nil, errors.New("password incorrect")
}
