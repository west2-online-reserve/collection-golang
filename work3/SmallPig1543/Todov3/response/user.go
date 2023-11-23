package response

type UserResp struct {
	ID       uint   `json:"id" form:"id" example:"1"`   // 用户ID
	UserName string `json:"user_name" form:"user_name"` // 用户名
	CreateAt int64  `json:"create_at" form:"create_at"` // 创建
}

type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
