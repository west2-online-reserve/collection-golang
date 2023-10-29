package service

import (
	"Todov3/DataBase"
	"Todov3/model"
	"Todov3/response"
	"Todov3/types"
	"Todov3/util"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
}

func (service *UserService) Register(req *types.UserRequest) (interface{}, error) {
	user, err := DataBase.FindUserByUsername(req.UserName)
	if err == nil {
		return response.BadResponse("该用户已存在"), err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.User{UserName: req.UserName}
		if err := user.SetPassword(req.PassWord); err != nil {
			return response.BadResponse("密码不符合要求"), err
		}
		if err := DataBase.CreateUser(user); err != nil {
			return response.BadResponse("存入数据库失败"), err
		}
	}
	return response.SuccessResponse(), nil
}

func (service *UserService) Login(req *types.UserRequest) (interface{}, error) {
	user, err := DataBase.FindUserByUsername(req.UserName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.BadResponse("该用户不存在"), err
	}
	if !user.CheckPassword(req.PassWord) {
		err = errors.New("密码错误")
		return response.BadResponse("密码错误"), err
	}
	//签发token
	token, err := util.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return response.BadResponse("token签发失败"), err
	}
	userResp := &response.TokenData{
		User: &response.UserResp{
			ID:       user.ID,
			UserName: user.UserName,
			CreateAt: user.CreatedAt.Unix(),
		},
		Token: token,
	}
	return response.SuccessResponseWithData(userResp), nil
}
