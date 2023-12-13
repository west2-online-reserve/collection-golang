package serializer

import "todolist/model"

type User struct {
	ID uint `json:"id" form:"id" example:"1"`
	UserName string `json:"user_name" example:"FanOne"`
	Status string `json:"status" form:"status"`
	CreatedAt int64 `json:"create_at" form:"create_at"`
}
//BuildUser序列化用户
func BuildUser(user model.User) User {
	return User {
		ID: user.ID,
		UserName: user.UserName,
		CreatedAt: user.CreatedAt.Unix(),
	}
}