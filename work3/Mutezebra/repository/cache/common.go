package cache

import (
	"github.com/go-redis/redis"
	"three/config"
	"three/pkg/utils"
)

var RedisClient *redis.Client

func InitRedis() {
	conf := config.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + conf.RedisPort,
		DB:       conf.RedisDbName,
		Network:  conf.RedisNetwork,
		Password: conf.RedisPassword,
	})
	_, err := client.Ping().Result()
	if err != nil {
		utils.LogrusObj.Fatalln(err)
	}
	RedisClient = client
	utils.LogrusObj.Infoln("Redis init success!")
}
