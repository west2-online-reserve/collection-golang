package pack

import (
	"Demo/biz/dal/db"
	"Demo/biz/model/model"
	"strconv"
)

func User(data *db.User) *model.User {
	return &model.User{
		ID:        data.Id,
		Username:  data.Username,
		Password:  data.Password,
		CreatedAt: strconv.FormatInt(data.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(data.UpdatedAt.Unix(), 10),
	}
}
