package rabbitmq

import (
	"context"
	"four/consts"
	"four/pkg/log"
	"four/repository/cache"
	"four/repository/db/dao"
	"github.com/streadway/amqp"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.LogrusObj.Errorf("%s:%s", err, msg)
	}
}

func PublishVideoInfoMsg(msg string) {
	ch, err := RabbitMQ.Channel()
	if err != nil {
		log.LogrusObj.Errorln(err)
	}
	err = ch.Publish("", consts.QueueNameOfCleanVideo, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		log.LogrusObj.Errorln(err)
	}
	ch.Close()
}

func CleanVideoInfo() {
	channel, err := RabbitMQ.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := channel.QueueDeclare(consts.QueueNameOfCleanVideo, true, false, false, false, nil)
	// 更改autoAck参数为false，关闭自动消息确认
	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		body := string(d.Body)
		form := body[0]
		body = body[1:]
		id, _ := strconv.ParseUint(body, 10, 32)
		vid := uint(id)
		if form == consts.DestroyCache {
			err = cache.DestroyVideoInfoCache(vid)
			if err != nil {
				log.LogrusObj.Errorln(err)
			}
		} else if form == consts.UpdateMysql {
			views, err := cache.RedisClient.Get(cache.VideoViewKey(vid)).Result()
			view, _ := strconv.Atoi(views)
			if err != nil {
				log.LogrusObj.Errorln(err)
			} else {
				videoDao := dao.NewVideoDao(context.Background())
				err = videoDao.UpdateVideoViews(vid, view)
				if err != nil {
					log.LogrusObj.Errorln(err)
				}
			}
		}
		d.Ack(false) // 手动传递消息确认
	}

	log.LogrusObj.Infoln("The clean goroutine over")
}
