// 用于定义配置相关的功能 
package conf

import (
    "fmt"
    "memo/model"
	"strings"
	"gopkg.in/ini.v1"
)

var (
	AppMode string
    HttpPort string
    Db string
    DbHost string
    DbPort string
    DbUser string
    DbPassWord string
    DbName string
)

//配置包的入口点
func Init() {
    //加载配置文件
    file,err:=ini.Load("./conf/config.ini")
    if err!=nil {
        fmt.Println("配置文件读取错误，请检查文件路径")
    }
    LoadServer(file)
    LoadMysql(file)
    //建数据库连接字符串path，包含数据库的用户、密码、地址、端口和名称等信息
    path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	fmt.Println(path)
    //传入数据库连接字符串,建立数据库连接 
	model.Database(path)
}

//加载服务器配置
func LoadServer(file *ini.File) {
    //从 [service] 部分读取 AppMode 和 HttpPort 的值
    AppMode=file.Section("service").Key("AppMode").String()
    HttpPort=file.Section("service").Key("HttpPort").String()
}

//加载MySQL配置 
func LoadMysql(file *ini.File) { 
    //从 [mysql] 部分读取数据库连接所需的各个参数
    Db=file.Section("mysql").Key("Db").String()
    DbHost=file.Section("mysql").Key("DbHost").String()
    DbPort=file.Section("mysql").Key("DbPort").String()
    DbUser=file.Section("mysql").Key("DbUser").String()
    DbPassWord=file.Section("mysql").Key("DbPassWord").String()
    DbName=file.Section("mysql").Key("DbName").String()
}