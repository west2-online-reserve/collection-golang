package rabbitmq

import (
	"four/config"
	"four/consts"
	"github.com/streadway/amqp"
	"log"
)

var RabbitMQ *amqp.Connection

func InitRabbitMQ() {
	// 1.尝试连接到RabbitMQ，建立连接
	// 该链接抽象了套接字连接，并为我们处理协议版本和认证等
	conf := config.Config.RabbitMQ
	url := conf.RabbitMQ + "://" + conf.RabbitMQUser + ":" + conf.RabbitMQPassWord + "@" + conf.RabbitMQHost + ":" + conf.RabbitMQPort + "/"
	conn, err := amqp.Dial(url)
	RabbitMQ = conn
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ,Errr:%s", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel,Errr:%s", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(consts.QueueNameOfCleanVideo, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue,Errr:%s", err)
	}
	for i := 0; i < consts.CleanVideoInfoGoroutineNumber; i++ {
		go CleanVideoInfo()
	}
	log.Println("Init Rabbitmq Success")
}
