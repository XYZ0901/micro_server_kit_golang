package initialize

import "go.uber.org/zap"

type Config struct {
	Name         string       `yaml:"name"`
	ServerConfig serverConfig `yaml:"server_config"`
}

type serverConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var (
	Logger *zap.Logger
	Cfg    Config
)

func init() {
	zapInit()
	cfgInit()
}
