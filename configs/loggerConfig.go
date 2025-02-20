package configs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level      string
	Encoding   string
	OutputPath string
	ErrorPath  string
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      "debug",
		Encoding:   "json",
		OutputPath: "stdout",
		ErrorPath:  "stderr",
	}
}

func (c *LoggerConfig) BuildLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(c.Level)); err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:             level,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          c.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{c.OutputPath},
		ErrorOutputPaths:  []string{c.ErrorPath},
	}

	return config.Build()
}
