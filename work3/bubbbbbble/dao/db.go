package dao

import (
	"bubbbbbble/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
type Info struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}
var DB *gorm.DB

func InitMysql() error{
	info := Info{
	Username: config.Vp.GetString("database.username"),
	Password: config.Vp.GetString("database.password"),
	Ip: config.Vp.GetString("database.ip"),
	Port: config.Vp.GetString("database.port"),
	Dbname: config.Vp.GetString("database.dbname"),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", info.Username, info.Password, info.Ip, info.Port, info.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		return err
	}
	DB = db
	return nil
}
