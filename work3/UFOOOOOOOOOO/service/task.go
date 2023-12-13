package service

import (
	"todolist/model"
	"todolist/serializer"
	"time"
)

type CreateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"`
}

type ShowTaskService struct {

}

type ListTaskService struct {
	PageNum int `json:"page_num" form:"page_name"`
	PageSize int `json:"page_size" form:"page_size"`
}


type UpdateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"`//0是未做，1是已做
}

type SearchTaskService struct {
	Info string `json:"info" form:"info"`
	PageNum int `json:"page_num" form:"page_name"`
	PageSize int `json:"page_size" form:"page_size"`
}


type DeleteTaskService struct {

}

//新增一条备忘录
func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	code := 200
	model.DB.First(&user, id)
	task := model.Task{
		User : user,
		Uid : user.ID,
		Title : service.Title,
		Status: 0,
		Content: service.Content,
		StartTime: time.Now().Unix(),
		EndTime: 0,
	}
	err := model.DB.Create(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg: "创建备忘录失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg: "创建成功",
	}
}

//展示一条备忘录
func (service *ShowTaskService) Show(tid string) serializer.Response{
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg: "查询失败",
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.BuildTask(task),
	}
}

//列表返回用户所有备忘录
func (service *ListTaskService) List(uid uint) serializer.Response{
	var tasks []model.Task
	count := 0
	if service.PageSize == 0{
		service.PageSize = 15
	}
	model.DB.Model(&model.Task{}).Preload("User").Where("uid=?",uid).
	Count(&count).Limit(service.PageSize).Offset((service.PageNum-1)*service.PageSize).Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

//更新备忘录
func (service *UpdateTaskService) Update(tid string) serializer.Response{
	var task model.Task
	model.DB.First(&task, tid)
	task.Content = service.Content
	task.Title = service.Title
	task.Status = service.Status
	model.DB.Save(&task)
	return serializer.Response{
		Status: 200,
		Data: serializer.BuildTask(task),
		Msg: "更新完成",
	}
}

//查询备忘录操作
func (service *SearchTaskService) Search(uid uint) serializer.Response{
	var tasks []model.Task
	count := 0
	if service.PageSize == 0 {
		service.PageSize = 10
	}
	model.DB.Model(&model.Task{}).Preload("User").Where("uid=?",uid).
	Where("title LIKE ? OR content LIKE?", "%"+service.Info+"%", 
	"%"+service.Info+"%").Count(&count).Limit(service.PageSize).Offset((service.PageNum-1)*service.PageSize).Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

//删除备忘录
func (service *DeleteTaskService) Delete(tid string) serializer.Response{
	var task model.Task
	err := model.DB.Delete(&task, tid).Error
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg: "删除失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg: "删除成功",
	}
}