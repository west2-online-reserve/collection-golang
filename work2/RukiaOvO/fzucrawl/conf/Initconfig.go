package conf

import (
	"github.com/spf13/viper"
)

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func InitConfig() *MysqlConfig {
	var one MysqlConfig

	dbSettings := viper.New()
	dbSettings.AddConfigPath("./conf/settings")
	dbSettings.SetConfigName("config")
	dbSettings.SetConfigType("yaml")

	if err := dbSettings.ReadInConfig(); err != nil {
		panic(err)
	}
	one.Username = dbSettings.GetString("DB_config.username")
	one.Password = dbSettings.GetString("DB_config.password")
	one.Host = dbSettings.GetString("DB_config.host")
	one.Port = dbSettings.GetInt("DB_config.port")
	one.Database = dbSettings.GetString("DB_config.database")

	return &one
}
