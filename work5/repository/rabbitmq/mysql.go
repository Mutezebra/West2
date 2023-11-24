package rabbitmq

import (
	"encoding/json"
	"five/consts"
	"five/pkg/log"
	"five/repository/db/dao"
	"five/repository/db/model"
	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/streadway/amqp"
	"strconv"
)

var PublishChannel []*amqp.Channel
var readMysqlQueueName []string
var consumerNumber int

// aReadMysqlQueueName 随机获取一个读取mysql的消费者
func aReadMysqlQueueName() string {
	return readMysqlQueueName[fastrand.Intn(consumerNumber)]
}

// PublishMessage 发布消息到rabbitmq
func PublishMessage(msg *model.Message) {
	ch := PublishChannel[fastrand.Intn(consumerNumber)]
	body, err := json.Marshal(msg)
	if err != nil {
		log.LogrusObj.Errorln(err)
		return
	}
	err = ch.Publish("", aReadMysqlQueueName(), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		log.LogrusObj.Errorln(err)
	}
}

// InitAsyncWriteMysql 异步写入mysql数据, count为启动的goroutine数量
func InitAsyncWriteMysql(count int) {
	consumerNumber = count

	// 初始化读取mysql的consumer名称，减少重复拼写字符串的开销
	for i := 0; i < count; i++ {
		readMysqlQueueName = append(readMysqlQueueName, consts.WriteMysqlQueue+"_"+strconv.Itoa(i))
	}

	// 初始化写入mysql的goroutine
	for i := 0; i < count; i++ {
		go AsyncWriteMysql(i)
	}

	// 初始化写入mysql的channel
	for i := 0; i < count; i++ {
		ch, err := RabbitMQConn.Channel()
		if err != nil {
			log.LogrusObj.Errorln(err)
			return
		}
		PublishChannel = append(PublishChannel, ch)
	}
}

// AsyncWriteMysql 异步写入mysql数据, i为goroutine的编号
func AsyncWriteMysql(i int) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.LogrusObj.Errorln(err)
		return
	}

	q, err := ch.QueueDeclare(
		consts.WriteMysqlQueue+"_"+strconv.Itoa(i),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.LogrusObj.Errorln(err)
		return
	}

	consumer, err := ch.Consume(
		q.Name,
		consts.ReadMysqlConsumer+"_"+strconv.Itoa(i),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.LogrusObj.Errorln(err)
		return
	}
	for {
		select {
		case msg := <-consumer:
			info := model.Message{}
			if err = json.Unmarshal(msg.Body, &info); err != nil {
				log.LogrusObj.Errorln(err)
				if err = msg.Ack(false); err != nil {
					log.LogrusObj.Errorln(err)
				}
				continue
			}
			if err = dao.SaveSingleChatMessage(&info); err != nil {
				log.LogrusObj.Errorln(err)
				continue
			}
			if err = msg.Ack(false); err != nil {
				log.LogrusObj.Errorln(err)
			}
		}
	}
}
