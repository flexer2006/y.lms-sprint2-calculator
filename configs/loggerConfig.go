// Package configs provides configuration structures and functions for the application.
package configs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig holds the configuration for the logger.
type LoggerConfig struct {
	Level      string // Level defines the logging level (e.g., debug, info, warn, error).
	Encoding   string // Encoding specifies the format of the log output (e.g., json, console).
	OutputPath string // OutputPath is the path where logs will be written.
	ErrorPath  string // ErrorPath is the path where error logs will be written.
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
