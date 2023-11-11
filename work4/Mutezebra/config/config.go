package config

import (
	"os"
)
import "github.com/spf13/viper"

var Config *Conf
var WorkDir string

type Conf struct {
	System   *System           `yaml:"system"`
	MySql    map[string]*MySql `yaml:"mysql"`
	Email    map[string]*Email `yaml:"email"`
	Redis    *Redis            `yaml:"redis"`
	Local    *Local            `yaml:"local"`
	ES       *ES               `yaml:"es"`
	QiNiu    *QiNiu            `yaml:"qiniu"`
	RabbitMQ *RabbitMQ         `yaml:"rabbitmq"`
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
	Domain    string `yaml:"domain"`
	HttpPort  string `yaml:"httpPort"`
	Host      string `yaml:"host"`
	LocalMode string `yaml:"LocalMode"`
}

type Email struct {
	Sender   string `yaml:"sender"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type Local struct {
	AvatarPath        string `yaml:"AvatarPath"`
	DefaultAvatarPath string `yaml:"DefaultAvatarPath"`
	DefaultVideoPath  string `yaml:"DefaultVideosPath"`
	QRCodePath        string `yaml:"QRCodePath"`
}

type ES struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type QiNiu struct {
	AccessKey         string `yaml:"AccessKey"`
	SecretKey         string `yaml:"SecretKey"`
	AvatarPath        string `yaml:"AvatarPath"`
	DefaultAvatarPath string `yaml:"DefaultAvatarPath"`
	VideoPath         string `yaml:"VideosPath"`
	Bucket            string `yaml:"Bucket"`
	Domain            string `yaml:"Domain"`
}

type RabbitMQ struct {
	RabbitMQ         string `yaml:"RabbitMQ"`
	RabbitMQUser     string `yaml:"RabbitMQUser"`
	RabbitMQPassWord string `yaml:"RabbitMQPassWord"`
	RabbitMQHost     string `yaml:"RabbitMQHost"`
	RabbitMQPort     string `yaml:"RabbitMQPort"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	WorkDir = workDir
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config/local")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}

func DirExistAndCreate(path string) error {
	dirPath := WorkDir + path
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0777)
	}
	return err
}

func DirInit() {
	err := DirExistAndCreate("/static")
	if err != nil {
		panic(err)
	}
	err = DirExistAndCreate("/static/imgs")
	if err != nil {
		panic(err)
	}
	err = DirExistAndCreate(Config.Local.AvatarPath)
	if err != nil {
		panic(err)
	}
	err = DirExistAndCreate(Config.Local.DefaultVideoPath)
	if err != nil {
		panic(err)
	}
	err = DirExistAndCreate(Config.Local.QRCodePath)
	if err != nil {
		panic(err)
	}
}
