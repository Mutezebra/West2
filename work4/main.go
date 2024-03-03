package main

import (
	"four/config"
	"four/pkg/log"
	"four/pkg/myutils"
	"four/pkg/regular"
	"four/repository/cache"
	"four/repository/db/dao"
	"four/repository/es"
	"four/repository/es/index"
	"four/repository/rabbitmq"
	"four/route"
)

func main() {
	initAll()
	r := route.NewRouter()
	r.Spin()
}

func initAll() {
	config.InitConfig()
	config.DirInit()
	log.InitLog()
	dao.InitMysql()
	cache.InitRedis()
	rabbitmq.InitRabbitMQ()
	regular.InitRegular()
	es.InitES()
	index.InitIndex()
	myutils.OssInit()
}
