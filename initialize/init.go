package initialize

import "go.uber.org/zap"

var (
	Logger *zap.Logger
)

func init() {
	zapInit()
}
