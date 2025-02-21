// Package logger provides a wrapper around zap.Logger with additional functionality.
package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger.
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
	opts  Options
}

var (
	globalLogger *Logger
	once         sync.Once
)

// Close shuts down the global logger and releases resources.
func Close() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Fatal logs a fatal message and writes it to a file before exiting.
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
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

	// Write the message and ensure it's synced to disk
	logger.Fatal(msg, fields...)
	if err := logger.Sync(); err != nil {
		l.Error("Failed to sync fatal log", zap.Error(err))
	}

	// This line will never be reached due to os.Exit in Fatal,
	// but we keep it as a fallback
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

	// Configure encoding
	config.Encoding = opts.Encoding
	config.OutputPaths = opts.OutputPath
	config.ErrorOutputPaths = opts.ErrorPath
	config.Development = opts.Development

	// Configure time format
	config.EncoderConfig = newEncoderConfig()

	// Create the logger
	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
		opts:   opts,
	}, nil
}

// GetLogger returns the global logger instance, creating it if necessary.
func GetLogger() *Logger {
	once.Do(func() {
		logger, err := New(DefaultOptions())
		if err != nil {
			fmt.Printf("Failed to create logger: %v\n", err)
			os.Exit(1)
		}
		globalLogger = logger
	})
	return globalLogger
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
