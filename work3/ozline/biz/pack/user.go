package pack

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/west2-online-reserve/collection-golang/work3/biz/dal/db"
	"github.com/west2-online-reserve/collection-golang/work3/biz/model/model"
)

func User(data *db.User) *model.User {
	hlog.Infof("data: %+v\n", data)
	return &model.User{
		ID:        data.Id,
		Username:  data.Username,
		Password:  data.Password,
		CreatedAt: strconv.FormatInt(data.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(data.UpdatedAt.Unix(), 10),
	}
}
