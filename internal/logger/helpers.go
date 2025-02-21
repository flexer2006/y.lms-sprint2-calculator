package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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