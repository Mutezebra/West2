package main

import (
	"five/config"
	"five/pkg/log"
	"five/repository/db/dao"
	"five/repository/rabbitmq"
	"five/repository/redis"
	"five/route"
	"five/service/ws"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	InitAll()
	r := route.NewRouter()
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	r.RunTLS(":8080", "./config/local/server.crt", "./config/local/server.key")
}

func InitAll() {
	config.InitConfig()
	log.InitLog()
	dao.InitMysql()
	redis.InitRedis()
	ws.InitWsService()
	rabbitmq.InitRabbitMQ()
	initTest()
}

func initTest() {
	ws.TestInit()
	rabbitmq.TestInit()
}
