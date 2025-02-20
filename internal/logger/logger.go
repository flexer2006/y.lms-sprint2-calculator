package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", "logs/error.log"}
	config.ErrorOutputPaths = []string{"stderr", "logs/error.log"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}

func Info(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Logger.Error(message, fields...)
}
