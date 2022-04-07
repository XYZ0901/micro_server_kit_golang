package initialize

import (
	"fmt"
	"github.com/streadway/amqp"
)

func rmqInit() {
	RmqConn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		Cfg.RMqConfig.User, Cfg.RMqConfig.Pwd, Cfg.RMqConfig.Host, Cfg.RMqConfig.Port))
	if err != nil {
		Logger.Sugar().Fatalf("Failed to connect rabbitMQ: %s\n", err.Error())
	}
}
