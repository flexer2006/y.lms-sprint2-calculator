package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// createNonExitingLogger creates a logger that doesn't exit on Fatal for testing
func createNonExitingLogger(t *testing.T, dir string) *logger.Logger {
	// Setup options
	opts := logger.Options{
		Level:      logger.Info,
		Encoding:   "json",
		OutputPath: []string{"stdout"},
		ErrorPath:  []string{"stderr"},
		LogDir:     dir,
	}

	// Create the logger
	l, err := logger.New(opts)
	require.NoError(t, err)

	// Replace the Fatal function with our test version
	l.Fatal = func(msg string, fields ...zapcore.Field) {
		// Log at error level instead of fatal to avoid os.Exit
		l.Error("Testing fatal: "+msg, fields...)

		// Create the log file manually to match the expected format
		if err := os.MkdirAll(dir, 0755); err == nil {
			timestamp := time.Now().Format("02-01-2006_15-04-05")
			logFile := filepath.Join(dir, fmt.Sprintf("fatal_%s.log", timestamp))

			fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			})

			file, err := os.Create(logFile)
			if err == nil {
				defer func() {
					if err := file.Close(); err != nil {
						l.Error("Failed to close log file", zap.Error(err))
					}
				}()

				// Create a simple entry that matches the expected format
				entry := zapcore.Entry{
					Level:   zapcore.FatalLevel,
					Time:    time.Now(),
					Message: msg,
				}

				// Write directly to file with error handling
				if buf, err := fileEncoder.EncodeEntry(entry, fields); err == nil {
					if _, err := file.Write(buf.Bytes()); err != nil {
						l.Error("Failed to write to log file", zap.Error(err))
					}
				} else {
					l.Error("Failed to encode log entry", zap.Error(err))
				}
			}
		}
	}

	return l
}
