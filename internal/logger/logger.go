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

// LogLevel определяет уровень логирования
type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

// Options определяет настройки логгера
type Options struct {
	Level       LogLevel
	Encoding    string
	OutputPath  []string
	ErrorPath   []string
	Development bool
	LogDir      string
}

// Logger представляет собой обертку над zap.Logger
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
	opts  Options
}

var (
	globalLogger *Logger
	once         sync.Once
)

// Close закрывает глобальный логгер и освобождает ресурсы
func Close() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

func DefaultOptions() Options {
	return Options{
		Level:       Info,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: false,
		LogDir:      "logs",
	}
}

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
	} else {
		defer file.Close()

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zapcore.FatalLevel,
		)

		combinedCore := zapcore.NewTee(l.Core(), fileCore)
		logger := zap.New(combinedCore)

		logger.Fatal(msg, fields...)
	}

	l.Logger.Fatal(msg, fields...)
}

// New создает новый логгер с указанными настройками
func New(opts Options) (*Logger, error) {
	config := zap.NewProductionConfig()

	// Устанавливаем уровень логирования
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

	// Настраиваем кодировку
	config.Encoding = opts.Encoding
	config.OutputPaths = opts.OutputPath
	config.ErrorOutputPaths = opts.ErrorPath
	config.Development = opts.Development

	// Настраиваем формат времени
	config.EncoderConfig = newEncoderConfig()

	// Создаем логгер
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

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

func (l *Logger) Sync() error {
	err1 := l.Logger.Sync()
	err2 := l.sugar.Sync()
	if err1 != nil {
		return err1
	}
	return err2
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("02-01-2006 15:04:05"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func extractContextFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field

	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields = append(fields, zap.Any("trace_id", traceID))
	}

	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, zap.Any("request_id", requestID))
	}

	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		fields = append(fields, zap.Any("correlation_id", correlationID))
	}

	return fields
}

// Close закрывает логгер и освобождает ресурсы
func (l *Logger) Close() error {
	return l.Sync()
}
