package config

import (
	"github.com/spf13/viper"
	"os"
	"three/pkg/utils"
)

var Config *Conf

type Conf struct {
	System *System           `yaml:"system"`
	MySql  map[string]*MySql `yaml:"mysql"`
	Redis  *Redis            `yaml:"redis"`
}

type MySql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

type System struct {
	Domain   string `yaml:"domain"`
	HttpPort string `yaml:"httpPort"`
	Host     string `yaml:"host"`
}

// InitConfig 读取配置文件到全局变量
func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wordDir + "/config/local")
	viper.AddConfigPath(wordDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
	utils.LogrusObj.Infoln("Config init success!")
}
