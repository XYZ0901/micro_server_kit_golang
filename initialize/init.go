package initialize

import (
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	Name         string       `yaml:"name"`
	ServerConfig serverConfig `yaml:"server_config"`
	MysqlConfig  mysqlConfig  `yaml:"mysql_config"`
	RedisConfig  redisConfig  `yaml:"redis_config"`
	RMqConfig    rmqConfig    `yaml:"rmq_config"`
}

type serverConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type mysqlConfig struct {
	User      string `yaml:"user"`
	Pwd       string `yaml:"pwd"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	DbName    string `yaml:"dbname"`
	Arguments string `yaml:"arguments"`
}

type redisConfig struct {
	Addr     string `yaml:"addr"`
	PassWord string `yaml:"pass_word"`
}

type rmqConfig struct {
	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	Logger      *zap.Logger
	Cfg         Config
	MysqlDb     *gorm.DB
	RedisClient *redis.Client
	RmqConn     *amqp.Connection
)

func init() {
	zapInit()
	cfgInit()
	sqlInit()
	redisInit()
	rmqInit()
}
