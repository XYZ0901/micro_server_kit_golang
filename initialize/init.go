package initialize

import (
	"github.com/go-redis/redis/v8"
	consul "github.com/hashicorp/consul/api"
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
	ConsulConfig consulConfig `yaml:"consul_config"`
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

type consulConfig struct {
	Addr    string `yaml:"addr"`
	Service struct {
		Tags []string `yaml:"tags"`
		Name string   `yaml:"name"`
	} `yaml:"service"`
}

type nacosConfig struct {
	Nacos struct {
		Server struct {
			Host        string `yaml:"host"`
			Port        int    `yaml:"port"`
			Scheme      string `yaml:"scheme"`
			ContextPath string `yaml:"context_path"`
		} `yaml:"server"`
		Client struct {
			NamespaceId         string `yaml:"namespace_id"`
			TimeoutMs           int    `yaml:"timeout_ms"`
			NotLoadCacheAtStart bool   `yaml:"not_load_cache_at_start"`
			LogDir              string `yaml:"log_dir"`
			CacheDir            string `yaml:"cache_dir"`
			LogLevel            string `yaml:"log_level"`
		} `yaml:"client"`
		Data struct {
			Group  string `yaml:"group"`
			DataId string `yaml:"data_id"`
			Type   string `yaml:"type"`
		}
	} `yaml:"nacos"`
}

var (
	Logger      *zap.Logger
	Cfg         Config
	MysqlDb     *gorm.DB
	RedisClient *redis.Client
	RmqConn     *amqp.Connection
	ConsulCli   *consul.Client
)

func init() {
	zapInit()
	cfgInit()
	sqlInit()
	redisInit()
	rmqInit()
	consulInit()
	sentinelInit()
}
