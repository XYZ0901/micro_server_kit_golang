package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	Name         string       `yaml:"name"`
	ServerConfig serverConfig `yaml:"server_config"`
	MysqlConfig  mysqlConfig  `yaml:"mysql_config"`
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

var (
	Logger  *zap.Logger
	Cfg     Config
	MysqlDb *gorm.DB
)

func init() {
	zapInit()
	cfgInit()
	sqlInit()
}
