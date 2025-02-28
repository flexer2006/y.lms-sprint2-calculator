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

func TestWithContext(t *testing.T) {
	t.Parallel()

	core, recorded := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)

	l, err := logger.New(logger.DefaultOptions())
	require.NoError(t, err)

	l.Logger = zapLogger

	ctx := context.WithValue(context.Background(), logger.RequestIDKey, "test-id")
	contextLogger := l.WithContext(ctx)

	contextLogger.Info("Test message")

	require.Equal(t, 1, recorded.Len())
	entry := recorded.All()[0]
	assert.Equal(t, "Test message", entry.Message)

	found := false
	for _, field := range entry.Context {
		if field.Key == "request_id" && field.String == "test-id" {
			found = true
			break
		}
	}
	assert.True(t, found, "Context fields were not properly added to log")
}

func TestSugaredLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)

	zapLogger := zap.New(core)

	sugar := zapLogger.Sugar()

	sugar.Infow("test message", "key", "value")

	if err := sugar.Sync(); err != nil {
		t.Fatalf("Failed to sync logger: %v", err)
	}

	output := buf.String()
	var logMap map[string]interface{}
	err := json.Unmarshal([]byte(output), &logMap)
	require.NoError(t, err)

	assert.Equal(t, "test message", logMap["msg"])
	assert.Equal(t, "value", logMap["key"])
}

func TestFatalLogging(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to clean up temporary directory: %v", err)
			// Don't fail the test on cleanup errors, but log them
		}
	}()

	testLog := createNonExitingLogger(t, tmpDir)

	testLog.Fatal("fatal test message", zap.String("test", "value"))

	time.Sleep(100 * time.Millisecond)

	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(files), 1)

	var fatalLogFile string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "fatal_") {
			fatalLogFile = filepath.Join(tmpDir, file.Name())
			break
		}
	}

	require.NotEmpty(t, fatalLogFile, "Fatal log file was not created")

	content, err := os.ReadFile(fatalLogFile)
	require.NoError(t, err)

	var logData map[string]interface{}
	err = json.Unmarshal(content, &logData)
	require.NoError(t, err)

	assert.Equal(t, "fatal test message", logData["msg"])
	assert.Equal(t, "value", logData["test"])
	assert.Equal(t, "fatal", logData["level"])
}
