package model

type User struct {
	ID       uint64 `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username;type:varchar(50);not null"`
	Password string `gorm:"column:password;type:varchar(100);not null"`
}

type Login struct {
	Username string `form:"username,required" json:"username,required"`
	Password string `form:"password,required" json:"password,required"`
}
