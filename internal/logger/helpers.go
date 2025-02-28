// Package logger provides logging utilities and helpers for the application.
package logger

import (
	"context"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is a custom type for context keys.
type ctxKey string

// Context key constants used for logging context fields.
const (
	TraceIDKey       ctxKey = "trace_id"
	RequestIDKey     ctxKey = "request_id"
	CorrelationIDKey ctxKey = "correlation_id"
)

// newEncoderConfig creates a new zapcore.EncoderConfig with custom settings.
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       common.LogFieldTimestamp,
		LevelKey:      common.LogFieldLevel,
		NameKey:       common.LogFieldLogger,
		CallerKey:     common.LogFieldCaller,
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    common.LogFieldMessage,
		StacktraceKey: common.LogFieldStacktrace,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("02-01-2006 15:04:05"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// extractContextFields extracts logging fields from the context.
func extractContextFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		fields = append(fields, zap.String(common.FieldTraceID, traceID.(string)))
	} else if traceID := ctx.Value(string(TraceIDKey)); traceID != nil {
		fields = append(fields, zap.String(common.FieldTraceID, traceID.(string)))
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields = append(fields, zap.String(common.FieldRequestID, requestID.(string)))
	} else if requestID := ctx.Value(string(RequestIDKey)); requestID != nil {
		fields = append(fields, zap.String(common.FieldRequestID, requestID.(string)))
	}

	if correlationID := ctx.Value(CorrelationIDKey); correlationID != nil {
		fields = append(fields, zap.String(common.FieldCorrelationID, correlationID.(string)))
	} else if correlationID := ctx.Value(string(CorrelationIDKey)); correlationID != nil {
		fields = append(fields, zap.String(common.FieldCorrelationID, correlationID.(string)))
	}

	return fields
}
