package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online-reserve/collection-golang/work3/biz/dal/db"
	"github.com/west2-online-reserve/collection-golang/work3/biz/model/user"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Login(req *user.LoginRequest) (*db.User, error) {
	return db.LoginCheck(s.ctx, req.Username, req.Password)
}

func (s *UserService) Register(req *user.RegisterRequest) (*db.User, error) {
	return db.CreateUser(s.ctx, req.Username, req.Password)
}
