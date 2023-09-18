package types

type UserRegisterReq struct {
	UserName string `json:"user_name" form:"user_name" binding:"required,min=1,max=10"`
	Password string `json:"password" form:"password" binding:"required,min=5,max=20"`
	NickName string `json:"nick_name" form:"nick_name" binding:"required,min=1,max=20"`
}

type UserLoginReq struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

type TokenDataResp struct {
	User         *UserResp `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserResp struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nick_name"`
	Password string `json:"password,omitempty"`
}
