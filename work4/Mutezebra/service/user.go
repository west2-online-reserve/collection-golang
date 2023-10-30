package service

import (
	"bytes"
	"context"
	"errors"
	"four/config"
	"four/consts"
	"four/pkg/ctl"
	"four/pkg/e"
	"four/pkg/myutils"
	"four/repository/db/dao"
	"four/repository/db/model"
	"four/repository/es/doc"
	esmodel "four/repository/es/model"
	"four/types"
	"gorm.io/gorm"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
)

type UserSrv struct {
}

var UserSrvOnce sync.Once

var UserService *UserSrv

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserService = &UserSrv{}
	})
	return UserService
}

func (s *UserSrv) Register(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	code := e.SUCCESS
	ok := myutils.IsEmail(req.Email)
	if !ok {
		code = e.InvalidEmailFormat
		err = errors.New("invalid email format")
		return ctl.RespError(code, err), err
	}
	user := &model.User{}
	userDao := dao.NewUserDao(ctx)
	_, err = userDao.FindUserByUserName(req.UserName)

	if err == gorm.ErrRecordNotFound {
		err = user.SetPassword(req.Password)
		if err != nil {
			code = e.SetPasswordFailed
			return ctl.RespError(code, err), err
		}
		user.Avatar = config.Config.Local.DefaultAvatarPath
		user.UserName = req.UserName
		user.NickName = req.NickName
		user.Email = req.Email
		err = userDao.Create(user)
		if err != nil {
			code = e.CreateUserFailed
			return ctl.RespError(code, err), err
		}

		// 向elasticsearch中添加数据
		u := esmodel.User{
			Uid:      user.ID,
			UserName: user.UserName,
		}
		u.CreateTime()
		err = doc.DocCreate(&u)
		if err != nil {
			code = e.CreateDocFailed
			return ctl.RespError(code, err), err
		}
		return ctl.RespSuccess(code), nil
	} else if err == nil {
		code = e.UserExists
		return ctl.RespSuccess(code), errors.New("user exists")
	} else {
		code = e.UserExists
		return ctl.RespError(code, err), errors.New("user exists")
	}
}

func (s *UserSrv) UserNameLogin(ctx context.Context, req *types.UserNameLoginReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(req.UserName)
	if err != nil {
		code = e.UserDoNotExist
		return ctl.RespError(code, err), err
	}
	if user.TotpEnable {
		ok := myutils.VerifyOtp(req.OTP, user.TotpSecret)
		if ok {
			err = user.CheckPassword(req.Password)
			if err != nil {
				code = e.CheckPasswordFailed
				return ctl.RespError(code, err), errors.New("check password failed")
			}
		} else {
			code = e.VerifyOtpFailed
			err = errors.New("verify otp failed")
			return ctl.RespError(code, err), err
		}
	}

	aToken, rToken, err := myutils.GenerateToken(user.UserName, user.ID)
	if err != nil {
		code = e.GenerateTokenFailed
		return ctl.RespError(code, err), err
	}

	u := &types.UserResp{
		ID:         user.ID,
		UserName:   user.UserName,
		NickName:   user.NickName,
		Email:      user.Email,
		VideoCount: model.VideoCount(user.UserName),
	}

	data := &types.TokenDataResp{
		User:         u,
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
	return ctl.RespSuccessWithData(code, data), err
}

func (s *UserSrv) EmailLogin(ctx context.Context, req *types.UserEmailLoginReq) (resp interface{}, err error) {
	code := e.SUCCESS
	code = e.SUCCESS
	ok := myutils.IsEmail(req.Email)
	if !ok {
		code = e.InvalidEmailFormat
		err = errors.New("invalid email format")
		return ctl.RespError(code, err), err
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserEmail(req.Email)
	if err != nil {
		code = e.UserDoNotExist
		return ctl.RespError(code, err), err
	}
	if user.TotpEnable {
		ok := myutils.VerifyOtp(req.OTP, user.TotpSecret)
		if ok {
			err = user.CheckPassword(req.Password)
			if err != nil {
				code = e.CheckPasswordFailed
				return ctl.RespError(code, err), errors.New("check password failed")
			}
		} else {
			code = e.VerifyOtpFailed
			err = errors.New("verify otp failed")
			return ctl.RespError(code, err), err
		}
	}

	aToken, rToken, err := myutils.GenerateToken(user.UserName, user.ID)
	if err != nil {
		code = e.GenerateTokenFailed
		return ctl.RespError(code, err), err
	}

	u := &types.UserResp{
		ID:         user.ID,
		UserName:   user.UserName,
		NickName:   user.NickName,
		Email:      user.Email,
		VideoCount: model.VideoCount(user.UserName),
	}

	data := &types.TokenDataResp{
		User:         u,
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
	return ctl.RespSuccessWithData(code, data), err
}

func (s *UserSrv) EnableTotp(ctx context.Context, req *types.UserEnableTotpReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}

	user, err := userDao.FindUserByUserID(userInfo.ID)
	if err != nil {
		code = e.FindUserFailed
		return ctl.RespError(code, err), err
	}

	if user.TotpEnable && req.Status == 1 { // 校验是否已经开启过
		code = e.ERROR
		err = errors.New("you have enabled")
		return ctl.RespError(code, err), err
	}

	status := true
	if req.Status == 0 { // 如果要关闭那么需要otp验证码
		ok := myutils.VerifyOtp(req.OTP, user.TotpSecret)
		if !ok {
			code = e.VerifyOtpFailed
			err = errors.New("verify otp failed")
			return ctl.RespError(code, err), err
		}
		status = false
	}

	err = userDao.UpdateTotpStatus(user, status) // 更新2FA状态
	if err != nil {
		code = e.UpdateTotpStatusFailed
		return ctl.RespError(code, err), err
	}

	if status { // 开启2FA
		key, err := myutils.GenerateOtp(user.UserName)
		if err != nil {
			code = e.GenerateOTPFailed
			return ctl.RespError(code, err), err
		}
		err = userDao.UpdateOtpSecret(user, key.Secret(), key.URL())
		if err != nil {
			code = e.UpdateOTPSecretFailed
			return ctl.RespError(code, err), err
		}
		img, _ := key.Image(200, 200)
		var buf bytes.Buffer
		_ = png.Encode(&buf, img)
		storePath := consts.OtpCodeStorePath + user.UserName + ".png"
		err = os.WriteFile(storePath, buf.Bytes(), 0644) // 把二维码写到文件里
		if err != nil {
			code = e.WriteFileFailed
			return ctl.RespError(code, err), err
		}
		err = myutils.SendEmail(user.Email, storePath) // 传入二维码的存储路径
		if err != nil {
			code = e.SendEmailFailed
			return ctl.RespError(code, err), err
		}
		return ctl.RespSuccessWithData(code, "Please wait you email"), err
	}
	return ctl.RespSuccessWithData(code, "已成功关闭"), nil
}

func (s *UserSrv) GetUserInfo(ctx context.Context, req *types.UserInfoReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	user, err := userDao.FindUserByUserName(userInfo.UserName)
	if err != nil {
		code = e.UserDoNotExist
		return ctl.RespError(code, err), err
	}

	data := types.UserResp{
		ID:         user.ID,
		UserName:   user.UserName,
		Avatar:     user.Avatar,
		Gender:     user.Gender,
		NickName:   user.NickName,
		Follow:     user.Follow,
		Fans:       user.Fans,
		Email:      user.Email,
		VideoCount: model.VideoCount(user.UserName),
	}
	return ctl.RespSuccessWithData(code, data), err
}

func (s *UserSrv) Update(ctx context.Context, req *types.UserUpdateReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	user, err := userDao.FindUserByUserName(userInfo.UserName)
	if err != nil {
		code = e.UserDoNotExist
		return ctl.RespError(code, err), err
	}
	switch req.Email {
	case "":
		break
	default:
		user.Email = req.Email
		break
	}
	switch req.Gender {
	case 0:
		break
	case 1:
		user.Gender = 1
		break
	case 2:
		user.Gender = 2
		break
	default:
		break
	}
	switch req.NickName {
	case "":
		break
	default:
		user.NickName = req.NickName
	}

	err = userDao.Update(user)
	if err != nil {
		code = e.UserInfoUpdateFailed
		return ctl.RespError(code, err), err
	}
	// 向elasticsearch中添加数据
	u := esmodel.User{
		Uid:      user.ID,
		UserName: user.UserName,
	}
	u.CreateTime()
	err = doc.DocCreate(&u)
	if err != nil {
		code = e.CreateDocFailed
		return ctl.RespError(code, err), err
	}

	data := types.UserResp{
		ID:         user.ID,
		UserName:   user.UserName,
		Avatar:     user.Avatar,
		Gender:     user.Gender,
		NickName:   user.NickName,
		Follow:     user.Follow,
		Fans:       user.Fans,
		Email:      user.Email,
		VideoCount: model.VideoCount(user.UserName),
	}

	return ctl.RespSuccessWithData(code, data), err
}

func (s *UserSrv) UpdateAvatar(ctx context.Context, avatar *multipart.FileHeader) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	user, err := userDao.FindUserByUserID(userInfo.ID)
	if err != nil {
		code = e.FindUserFailed
		return ctl.RespError(code, err), err
	}
	ext, err := myutils.ValidAvatar(avatar)
	if err != nil {
		code = e.UnValidAvatar
		return ctl.RespError(code, err), err
	}
	var path string
	if config.Config.System.LocalMode == "local" {
		path = config.Config.Local.AvatarPath + strconv.Itoa(int(userInfo.ID)) + ext
		err = myutils.SavedAvatarFile(avatar, path)
		if err != nil {
			code = e.SaveAvatarFailed
			return ctl.RespError(code, err), err
		}
	} else {
		name := strconv.Itoa(int(userInfo.ID)) + ext
		file, err := avatar.Open()
		if err != nil {
			code = e.AvatarFileOpenFailed
			return ctl.RespError(code, err), err
		}
		data, err := io.ReadAll(file)
		if err != nil {
			code = e.ReadAvatarFileFailed
			return ctl.RespError(code, err), err
		}
		_, err, path = myutils.UploadAvatar(ctx, name, data)
		if err != nil {
			code = e.OSSUploadVideoFailed
		}
	}
	err = userDao.UpdateAvatar(user, path)
	if err != nil {
		code = e.UpdateAvatarFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccess(code), err
}

func (s *UserSrv) Follow(ctx context.Context, req *types.UserFollowReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	userDao := dao.NewUserDao(ctx)
	fan := &model.Fans{Uid: userInfo.ID, FollowerId: req.FollowerID}
	err = userDao.Follow(fan)
	if err != nil {
		code = e.FollowFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccess(code), err
}

func (s *UserSrv) UnFollow(ctx context.Context, req *types.UserFollowReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	userDao := dao.NewUserDao(ctx)
	fan := &model.Fans{Uid: userInfo.ID, FollowerId: req.FollowerID}
	err = userDao.UnFollow(fan)
	if err != nil {
		code = e.UnFollowFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccess(code), err
}

func (s *UserSrv) GetFriendList(ctx context.Context, req *types.UserGetFriendReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	list, err := userDao.FriendList(userInfo.ID)
	if err != nil {
		code = e.GetFriendListFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccessWithData(code, list), err
}

func (s *UserSrv) GetFollowerList(ctx context.Context, req *types.UserGetFollowerReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	list, err := userDao.FollowerList(userInfo.ID)
	if err != nil {
		code = e.GetFollowerListFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccessWithData(code, list), err
}

func (s *UserSrv) Delete(ctx context.Context) (resp interface{}, err error) {
	code := e.SUCCESS
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}

	userDao := dao.NewUserDao(ctx)

	_, err = userDao.FindUserByUserID(userInfo.ID)
	if err != nil {
		code = e.UserDoNotExist
		return ctl.RespError(code, err), err
	}

	err = userDao.Delete(userInfo.ID)
	if err != nil {
		code = e.DeleteUserFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccess(code), nil
}
