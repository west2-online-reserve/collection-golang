package config

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
	// Log struct {
	// 	Level string `yaml:"level"`
	// 	Path  string `yaml:"path"`
	// } `yaml:"log"`
	Jwt struct {
		SecretKey  string        `yaml:"secretKey"`
		Signed     string        `yaml:"signed"`
		ExpireTime time.Duration `yaml:"expireTime"`
		MaxRefresh time.Duration `yaml:"maxRefresh"`
	} `yaml:"jwt"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
}

var (
	once     sync.Once
	instance *config
)

func GetConfig() *config {
	once.Do(func() {
		instance = &config{}
		if err := instance.load("config/config.yaml"); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	})
	return instance
}

func (c *config) load(path string) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
