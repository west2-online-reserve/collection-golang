package types

type UserResp struct {
	ID         uint   `json:"id,omitempty" `
	UserName   string `json:"user_name,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Gender     int    `json:"gender,omitempty"`
	NickName   string `json:"nick_name,omitempty"`
	Follow     int    `json:"follow,omitempty"`
	Fans       int    `json:"fans,omitempty"`
	Email      string `json:"email,omitempty"`
	VideoCount int    `json:"video_count"`
}

type UserInfoReq struct {
}

type TokenDataResp struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type UserRegisterReq struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password,required" form:"password"`
	NickName string `json:"nick_name,required" form:"nick_name"`
	Email    string `json:"email" form:"email,required"`
}

type UserNameLoginReq struct {
	UserName string `json:"user_name" form:"user_name,required"`
	Password string `json:"password" form:"password,required"`
	OTP      string `json:"otp" form:"otp"`
}

type UserEmailLoginReq struct {
	Email    string `json:"email" form:"email,required"`
	Password string `json:"password" form:"password,required"`
	OTP      string `json:"otp" form:"otp"`
}

type UserEnableTotpReq struct {
	// Status 0关闭,1开启
	Status int    `json:"status" form:"status"`
	OTP    string `json:"otp" form:"otp"`
}

type UserUpdateReq struct {
	Email    string `json:"email" form:"email"`
	NickName string `json:"nick_name" form:"nick_name"`
	Gender   int    `json:"gender" form:"gender"`
}

type UserAvatarUpdateReq struct {
}

type UserGetFriendReq struct {
}

type UserGetFollowerReq struct {
}

type UserGetFansReq struct {
}

type UserFollowReq struct {
	FollowerID uint `json:"follower_id" form:"follower_id"`
}
