//迁移到数据库里面
//确保数据库的结构与代码中定义的模型结构保持一致 
package model

//执行数据库迁移操作 
func migration() {
	//自动迁移模式 
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{}).
		AutoMigrate(&Task{})
	//外键关联 
	DB.Model(&Task{}).AddForeignKey("uid", "User(id)", "CASCADE", "CASCADE")
} 