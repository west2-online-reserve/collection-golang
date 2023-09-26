package service
import (
	"bubbbbbble/model"
	"bubbbbbble/dao"
)


var PageSize = 10

func Create(todo model.Todo) {
	dao.DB.Table("todos").Create(&todo)
}
func GetSingle(name string, id int) (todo model.Todo, err error) {
	err = dao.DB.Table("todos").Where("name = ? and id = ?", name, id).Find(&todo).Error
	todo.Viewed++
	dao.DB.Table("todos").Save(&todo)
	return todo, err
}
func GetAll(name string) (todolist []model.Todo, err error) {
	err = dao.DB.Table("todos").Where("name = ?", name).Limit(PageSize).Find(&todolist).Error
	if err != nil {
		return todolist, err
	}
	for _, todo := range todolist {
		todo.Viewed++
		dao.DB.Table("todos").Save(todo)
	}
	return todolist, err
}
func GetAllDone(name string) (todolist []model.Todo, err error) {
	err = dao.DB.Table("todos").Where("name = ? and status = 1", name).Limit(PageSize).Find(&todolist).Error
	if err != nil {
		return todolist, err
	}
	for _, todo := range todolist {
		todo.Viewed++
		dao.DB.Table("todos").Save(todo)
	}
	return todolist, err
}
func GetAllUndo(name string) (todolist []model.Todo, err error) {
	err = dao.DB.Table("todos").Where("name = ? and status = 0", name).Limit(PageSize).Find(&todolist).Error
	for _, todo := range todolist {
		todo.Viewed++
		dao.DB.Table("todos").Save(&todo)
	}
	return todolist, err
}
func GetByKey(name string, key string) (todolist []model.Todo, err error) {
	keystr := "%" + key + "%"
	err = dao.DB.Table("todos").Where("name = ? and (content like ? or title like ? or deadline like ? or starttime like ?)", name, keystr, keystr, keystr, keystr).Limit(PageSize).Find(&todolist).Error
	for _, todo := range todolist {
		todo.Viewed++
		dao.DB.Table("todos").Save(&todo)
	}
	return todolist, err
}
func UpdateAllDone(name string) {
	dao.DB.Table("todos").Where("name = ?", name).Update("status", 1)

}
func UpdateAllUndo(name string) {
	dao.DB.Table("todos").Where("name = ?", name).Update("status", 0)

}
func UpdateSingleDone(name string, id int) error {
	err := dao.DB.Table("todos").Where("name = ? and id = ? ", name, id).Update("status", 1).Error
	return err
}
func UpdateSingleUndo(name string, id int) error {
	err := dao.DB.Table("todos").Where("name = ? and id = ? ", name, id).Update("status", 0).Error
	return err
}
func DelAllDone(name string) {
	dao.DB.Table("todos").Where("name = ? and status = 1", name).Delete(&model.Todo{})
}
func DelAllUndo(name string) {
	dao.DB.Table("todos").Where("name = ? and status = 0", name).Delete(&model.Todo{})
}
func DelSingle(name string, id int) error {
	err := dao.DB.Table("todos").Where("name = ? and id = ?", name, id).Delete(&model.Todo{}).Error
	return err
}
func DelAll(name string) {
	dao.DB.Table("todos").Where("name = ?", name).Delete(&model.Todo{})
}

func Login(user *model.User) error{
	err := dao.DB.Table("users").Where("name = ? AND password = ?", user.Name, user.Password).First(&model.User{}).Error
	return err
}
func SignUp(user *model.User)error{
	err:=dao.DB.Table("users").Where("name = ?", user.Name).First(&model.User{}).Error
	return err
}
func CreateUser(user model.User){
	dao.DB.Table("users").Create(&user)
}