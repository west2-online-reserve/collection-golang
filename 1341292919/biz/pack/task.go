package pack

import (
	"Demo/biz/dal/db"
	"Demo/biz/model/model"
	"strconv"
)

func Task(data *db.Task) *model.Task {
	return &model.Task{
		ID:        data.Id,
		Title:     data.Title,
		Content:   data.Content,
		CreatedAt: strconv.FormatInt(data.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(data.UpdatedAt.Unix(), 10),
		StartAt:   strconv.FormatInt(data.StartAt.Unix(), 10),
		EndAt:     strconv.FormatInt(data.EndAt.Unix(), 10),
	}
}

func TaskList(data []*db.Task, total int64) *model.TaskList {
	resp := make([]*model.Task, 0, len(data))
	for _, v := range data {
		resp = append(resp, Task(v))
	}

	return &model.TaskList{
		Items: resp,
		Total: total,
	}
}
