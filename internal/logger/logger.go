// Package logger provides a wrapper around zap.Logger with additional functionality.
package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger.
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
	opts  Options
	Fatal func(msg string, fields ...zapcore.Field) // Add this field
}

// defaultFatal is the default implementation of Fatal
func (l *Logger) defaultFatal(msg string, fields ...zapcore.Field) {
	if err := os.MkdirAll(l.opts.LogDir, 0755); err != nil {
		l.Error("Failed to create logs directory", zap.Error(err))
	}

	timestamp := time.Now().Format("02-01-2006_15-04-05")
	logFile := filepath.Join(l.opts.LogDir, fmt.Sprintf("fatal_%s.log", timestamp))

	fileEncoder := zapcore.NewJSONEncoder(newEncoderConfig())
	file, err := os.Create(logFile)
	if err != nil {
		l.Error("Failed to create log file", zap.Error(err))
		l.Logger.Fatal(msg, fields...)
	}
	defer func() {
		if err := file.Close(); err != nil {
			l.Error("Failed to close log file", zap.Error(err))
		}
	}()

	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(file),
		zapcore.FatalLevel,
	)

	combinedCore := zapcore.NewTee(l.Core(), fileCore)
	logger := zap.New(combinedCore)

	logger.Fatal(msg, fields...)
	if err := logger.Sync(); err != nil {
		l.Error("Failed to sync fatal log", zap.Error(err))
	}

	l.Logger.Fatal(msg, fields...)
}

// New creates a new logger with the specified options.
func New(opts Options) (*Logger, error) {
	config := zap.NewProductionConfig()

	// Set the logging level
	switch opts.Level {
	case Debug:
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case Info:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case Warn:
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case Error:
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown log level: %s", opts.Level)
	}

	config.Encoding = opts.Encoding
	config.OutputPaths = opts.OutputPath
	config.ErrorOutputPaths = opts.ErrorPath
	config.Development = opts.Development

	config.EncoderConfig = newEncoderConfig()

	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	l := &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
		opts:   opts,
	}

	l.Fatal = l.defaultFatal

	return l, nil
}

// WithContext returns a new Logger with fields extracted from the context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := extractContextFields(ctx)
	if len(fields) == 0 {
		return l
	}

	newLogger := l.Logger.With(fields...)
	return &Logger{
		Logger: newLogger,
		sugar:  newLogger.Sugar(),
		opts:   l.opts,
	}
}

// Sugar returns the SugaredLogger for structured logging.
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	err1 := l.Logger.Sync()
	err2 := l.sugar.Sync()
	if err1 != nil {
		return err1
	}
	return err2
}

// Close shuts down the logger and releases resources.
func (l *Logger) Close() error {
	return l.Sync()
}
