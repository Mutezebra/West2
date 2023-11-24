package rabbitmq

import (
	"five/config"
	"five/consts"
	"five/pkg/log"
	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection

func InitRabbitMQ() {
	// 1.尝试连接到RabbitMQ，建立连接
	// 该链接抽象了套接字连接，并为我们处理协议版本和认证等
	conf := config.Config.RabbitMQ
	url := conf.RabbitMQ + "://" + conf.RabbitMQUser + ":" + conf.RabbitMQPassWord + "@" + conf.RabbitMQHost + ":" + conf.RabbitMQPort + "/"
	conn, err := amqp.Dial(url)
	RabbitMQConn = conn
	if err != nil {
		log.LogrusObj.Panic("Failed to connect to RabbitMQConn,Err:%s", err)
	}
	InitAsyncWriteMysql(consts.ASyncWriteMysqlNumber)
	log.LogrusObj.Infoln("Init Rabbitmq Success")
}

func TestInit() {
	go DoneSinglePublish()
}
