package regular

import (
	"four/consts"
	"four/pkg/log"
	"four/repository/cache"
	"four/repository/rabbitmq"
	"strconv"
	"strings"
	"time"
)

func CleanVideoRegular() {
	for {
		// 每60秒清理一次
		time.Sleep(consts.RegularClean)
		now := time.Now().Unix()
		// 获取所有的key
		keys, err := cache.RedisClient.Keys("info:video:*").Result()
		if err != nil {
			log.LogrusObj.Errorln("Error getting keys:", err)
			continue
		}

		// 遍历每一个 key 进行清理或者更新
		for _, key := range keys {
			split := strings.Split(key, ":")
			vid := split[2]
			result := cache.RedisClient.HGetAll(key)
			value, err := result.Result()
			if err != nil {
				log.LogrusObj.Errorln("Error getting keys:", err)
			}
			lastUse, _ := strconv.ParseInt(value["last_use"], 10, 64)
			update, _ := strconv.ParseInt(value["update"], 10, 64)
			if now-update >= 15 {
				msg := "1" + vid
				rabbitmq.PublishVideoInfoMsg(msg)
				cache.RedisClient.HSet(key, "update", time.Now().Unix()) // 刷新更新时间
			}
			if now-lastUse >= 60 {
				msg := "0" + vid
				rabbitmq.PublishVideoInfoMsg(msg)
			}
		}
	}
}
