package redis

import (
	"five/config"
	"five/pkg/log"
	"github.com/go-redis/redis"
	log2 "log"
)

var RedisClient *redis.Client

// InitRedis 初始化redis
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
	SomeValueInit()
	log.LogrusObj.Infoln("Redis init success!")
}

func SomeValueInit() {
	// 初始化用户id
	RedisClient.SetNX("user_id", 0, 0)
}
