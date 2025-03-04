// Package configs предоставляет конфигурационные структуры и функции для приложения.
package configs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig содержит конфигурацию для регистратора.
type LoggerConfig struct {
	Level      string // Level определяет уровень логирования(например, debug, info, warn, error).
	Encoding   string // Encoding задает формат вывода журнала (например, json, console).
	OutputPath string // OutputPath - это путь, по которому будут записываться журналы.
	ErrorPath  string // ErrorPath - путь, куда будут записываться журналы ошибок.
}

// BuildLogger creates and returns a new zap.Logger based on the LoggerConfig settings.
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
