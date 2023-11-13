package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"lyt112/FZUSpider/model"
)

var (
	Host     string
	Port     string
	UserName string
	Password string
	Database string
	bvNumber string
)

func InitDataBase() {
	file, err := ini.Load("./config/conf.ini")
	if err != nil {
		fmt.Println("获取数据库配置失败")
		return
	}
	LoadMysql(file)
}

func LoadMysql(file *ini.File) {
	Host = file.Section("mysql").Key("host").String()
	Port = file.Section("mysql").Key("port").String()
	UserName = file.Section("mysql").Key("username").String()
	Password = file.Section("mysql").Key("password").String()
	Database = file.Section("mysql").Key("database").String()
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", UserName, Password, Host, Port, Database)
	model.InitDB(dataSourceName)
}

func GetBvNumber() string {
	file, err := ini.Load("./config/conf.ini")
	if err != nil {
		fmt.Println("获取bv号失败")
		return ""
	}
	bvNumber = file.Section("Bilibili").Key("bvNumber").String()
	return bvNumber
}
