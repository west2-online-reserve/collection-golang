package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sync"
	"three/pkg/ctl"
	"three/pkg/e"
	"three/pkg/utils"
	"three/repository/db/dao"
	"three/types"
)

var userServiceOnce sync.Once
var UserSrvIns *UserService

type UserService struct {
}

func GetUserSrv() *UserService {
	userServiceOnce.Do(func() {
		UserSrvIns = &UserService{}
	})
	return UserSrvIns
}

func (s *UserService) Register(ctx context.Context, req *types.UserRegisterReq) (interface{}, error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(req.UserName)
	if err == gorm.ErrRecordNotFound {
		err = user.SetPassword(req.Password)
		if err != nil {
			code = e.SetPasswordFailed
			return ctl.RespError(err, code), err
		}
		user.UserName = req.UserName
		user.NickName = req.NickName
		err = userDao.Create(user)
		if err != nil {
			code = e.CreateUserFailed
			return ctl.RespError(err, code), err
		}
		return ctl.RespSuccess(code), nil
	} else if err == nil {
		code = e.UserExists
		return ctl.RespSuccess(code), errors.New("user exists")
	} else {
		code = e.UserExists
		return ctl.RespError(err, code), errors.New("user exists")
	}
}

func (s *UserService) Login(ctx context.Context, req *types.UserLoginReq) (interface{}, error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(req.UserName)
	if err != nil {
		code = e.UserNotExists
		return ctl.RespError(err, code), err
	}
	err = user.CheckPassword(req.Password)
	if err != nil {
		code = e.CheckPasswordFailed
		return ctl.RespError(err, code), err
	}
	accessToken, refreshToken, err := utils.GenerateToken(user.UserName, user.ID)
	if err != nil {
		code = e.GenerateTokenFailed
		return ctl.RespError(err, code), err
	}
	u := &types.UserResp{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
	}
	data := types.TokenDataResp{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctl.RespSuccessWithData(data, code), nil
}
