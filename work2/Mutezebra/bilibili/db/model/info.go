package model

type Info struct {
	Message string `gorm:"type:longtext"`
	Mid     int64
}
