package model
type Todo struct {
	Id        int    `gorm:"id" form:"id" json:"id"`
	Name      string `gorm:"name" form:"name" json:"name" binding:"required"`
	Title     string `gorm:"title" form:"title" json:"title" binding:"required"`
	Content   string `gorm:"content "  form:"content" json:"content" binding:"required"`
	Starttime string `gorm:"starttime "  form:"starttime" json:"starttime" binding:"required"`
	Deadline  string `gorm:"deadline"  form:"deadline" json:"deadline" binding:"required"`
	Viewed    int    `gorm:"viewed"  form:"viewed" json:"viewed"`
	Status    int    `gorm:"status"  form:"status" json:"status"  binding:"required"`
}
type User struct {
	Name     string `json:"name" form:"name" gorm:"name" binding:"required"`
	Password string `json:"password" form:"password" gorm:"password" binding:"required"`
}
