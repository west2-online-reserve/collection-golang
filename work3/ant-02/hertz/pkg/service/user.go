package service

import (
	"hertz/pkg/model"
	"hertz/pkg/repository"
	"hertz/util"
)

type userService struct {
	ur repository.UserRepository
}

type UserService interface {
	Login(username, password string) uint64
	Register(username, password string) bool
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{ur: ur}
}

func (us *userService) Login(username, password string) uint64 {
	u := &model.User{
		Username: username,
		Password: util.MD5Sum(password),
	}

	if us.ur.QueryUser(u) {
		return u.ID
	}
	return 0
}

func (us *userService) Register(username, password string) bool {
	u := &model.User{
		Username: username,
		Password: util.MD5Sum(password),
	}

	return us.ur.Register(u)
}
