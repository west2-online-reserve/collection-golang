package db

import (
	"Demo/pkg/constants"
	"context"
	"errors"
)

func CreateUser(ctx context.Context, username, password string) (*User, error) {
	var userResp *User

	err := DB.WithContext(ctx).
		Table(constants.TableUser).
		Where("BINARY username = ?", username).
		First(&userResp).
		Error

	if err == nil {
		return nil, errors.New("username already exists")
	}

	userResp = &User{
		Username: username,
		Password: password,
	}

	err = DB.WithContext(ctx).
		Table(constants.TableUser).
		Create(&userResp).
		Error
	if err != nil {
		return nil, err
	}

	return userResp, nil
}

func LoginCheck(ctx context.Context, username string, password string) (*User, error) {
	var userResp *User

	err := DB.WithContext(ctx).
		Table(constants.TableUser).
		Where("BINARY username = ?", username). //可以强制区分大小写
		First(&userResp).
		Error

	if err != nil {
		return nil, errors.New("username didn't exist")
	}

	if userResp.Password != password {
		return nil, errors.New("password didn't match")
	} else {
		return userResp, nil
	}
}
