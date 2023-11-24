package config

import "os"
import "github.com/spf13/viper"

var Config *Conf

type Conf struct {
	System   *System           `yaml:"system"`
	MySql    map[string]*MySql `yaml:"mysql"`
	Redis    *Redis            `yaml:"redis"`
	Local    *Local            `yaml:"local"`
	RabbitMQ *RabbitMQ         `yaml:"rabbitmq"`
}

type MySql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

type System struct {
	Domain    string `yaml:"domain"`
	HttpPort  string `yaml:"httpPort"`
	Host      string `yaml:"host"`
	LocalMode string `yaml:"LocalMode"`
}

type RabbitMQ struct {
	RabbitMQ         string `yaml:"RabbitMQ"`
	RabbitMQUser     string `yaml:"RabbitMQUser"`
	RabbitMQPassWord string `yaml:"RabbitMQPassWord"`
	RabbitMQHost     string `yaml:"RabbitMQHost"`
	RabbitMQPort     string `yaml:"RabbitMQPort"`
}

type Local struct {
	AvatarPath        string `yaml:"AvatarPath"`
	DefaultAvatarPath string `yaml:"DefaultAvatarPath"`
	DefaultVideoPath  string `yaml:"DefaultVideosPath"`
}

func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wordDir + "/config/local")
	viper.AddConfigPath(wordDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
