package service

import (
	db2 "Demo/biz/dal/db"
	"Demo/biz/model/user"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx, c}
}

func (s *UserService) Login(req *user.LoginRequest) (*db2.User, error) {
	return db2.LoginCheck(s.ctx, req.Username, req.Password)
}

func (s *UserService) Register(req *user.RegisterRequest) (*db2.User, error) {
	return db2.CreateUser(s.ctx, req.Username, req.Password)
}
