package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is a custom type for context keys
type ctxKey string

// Context key constants
const (
	TraceIDKey       ctxKey = "trace_id"
	RequestIDKey     ctxKey = "request_id"
	CorrelationIDKey ctxKey = "correlation_id"
)

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

	// Extract context values using both string and custom key types
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		fields = append(fields, zap.String("trace_id", traceID.(string)))
	} else if traceID := ctx.Value(string(TraceIDKey)); traceID != nil {
		fields = append(fields, zap.String("trace_id", traceID.(string)))
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields = append(fields, zap.String("request_id", requestID.(string)))
	} else if requestID := ctx.Value(string(RequestIDKey)); requestID != nil {
		fields = append(fields, zap.String("request_id", requestID.(string)))
	}

	if correlationID := ctx.Value(CorrelationIDKey); correlationID != nil {
		fields = append(fields, zap.String("correlation_id", correlationID.(string)))
	} else if correlationID := ctx.Value(string(CorrelationIDKey)); correlationID != nil {
		fields = append(fields, zap.String("correlation_id", correlationID.(string)))
	}

	return fields
}
