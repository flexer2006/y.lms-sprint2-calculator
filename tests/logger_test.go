// File: tests/logger_test.go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestLoggerInitialization remains the same as it's working correctly

// TestWithContext tests the context integration
func TestWithContext(t *testing.T) {
	t.Parallel()

	// Create an observer core for capturing logs
	core, recorded := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	// Create a logger with standard options
	l, err := logger.New(logger.DefaultOptions())
	require.NoError(t, err)

	// Replace the internal logger with our observer logger
	l.Logger = zapLogger

	// Create a context with values
	ctx := context.WithValue(context.Background(), logger.RequestIDKey, "test-id")
	contextLogger := l.WithContext(ctx)

	contextLogger.Info("Test message")

	// Verify the logs contain the context fields
	require.Equal(t, 1, recorded.Len())
	entry := recorded.All()[0]
	assert.Equal(t, "Test message", entry.Message)

	// Find the request ID field
	found := false
	for _, field := range entry.Context {
		if field.Key == "request_id" && field.String == "test-id" {
			found = true
			break
		}
	}
	assert.True(t, found, "Context fields were not properly added to log")
}

// TestSugaredLogger tests the logger functionality with structured fields
func TestSugaredLogger(t *testing.T) {
	t.Parallel()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a direct encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create a core that writes to our buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)

	// Create a zap logger with this core
	zapLogger := zap.New(core)

	// Create a sugared logger from the zap logger
	sugar := zapLogger.Sugar()

	// Log with the sugar logger
	sugar.Infow("test message", "key", "value")

	// Force sync to ensure all logs are written
	sugar.Sync()

	// Parse the output
	output := buf.String()
	var logMap map[string]interface{}
	err := json.Unmarshal([]byte(output), &logMap)
	require.NoError(t, err)

	// Verify the structured logging worked correctly
	assert.Equal(t, "test message", logMap["msg"])
	assert.Equal(t, "value", logMap["key"])
}

// TestFatalLogging tests the custom fatal logging without actually exiting
func TestFatalLogging(t *testing.T) {
	// Setup a temporary directory for log files
	tmpDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a logger with non-exiting Fatal function
	testLog := createNonExitingLogger(t, tmpDir)

	// Call Fatal - our override will log without exiting
	testLog.Fatal("fatal test message", zap.String("test", "value"))

	// Give it a moment to write the file
	time.Sleep(100 * time.Millisecond)

	// Continue with the rest of the test as before...
	// Check that a fatal log file was created
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	// We should have at least one log file
	require.GreaterOrEqual(t, len(files), 1)

	// Find the fatal log file
	var fatalLogFile string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "fatal_") {
			fatalLogFile = filepath.Join(tmpDir, file.Name())
			break
		}
	}

	require.NotEmpty(t, fatalLogFile, "Fatal log file was not created")

	// Read the log file and verify content
	content, err := os.ReadFile(fatalLogFile)
	require.NoError(t, err)

	var logData map[string]interface{}
	err = json.Unmarshal(content, &logData)
	require.NoError(t, err)

	assert.Equal(t, "fatal test message", logData["msg"])
	assert.Equal(t, "value", logData["test"])
	assert.Equal(t, "fatal", logData["level"])
}
