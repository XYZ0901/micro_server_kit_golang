package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func zapInit() {
	encode := getEncode()
	writerSyncer := getLogWrite()
	//TODO: this 'DebugLevel' you need to rewriter if you want
	core := zapcore.NewCore(encode, writerSyncer, zapcore.DebugLevel)
	Logger = zap.New(core, zap.AddCaller())
}

//TODO: this 'Encode' you need to rewriter if you want
func getEncode() zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderCfg)
}

// write log to ...
func getLogWrite() zapcore.WriteSyncer {
	//TODO: this File Path you need to change if you want
	file, err := os.OpenFile("./log/dev.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln("[init.getLogWrite] OpenFile ./log/dev.log error: ", err)
	}
	return zapcore.AddSync(file)
}
