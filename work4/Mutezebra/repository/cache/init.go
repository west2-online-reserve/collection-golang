package cache

import (
	"four/config"
	"four/pkg/log"
	"github.com/go-redis/redis"
	log2 "log"
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
		log2.Println(err)
		log.LogrusObj.Panic(err)
	}
	RedisClient = client
	log.LogrusObj.Infoln("Redis init success!")
}
