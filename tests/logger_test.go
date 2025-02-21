package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define custom key types for context values
type ctxKey string

const (
	traceIDKey   ctxKey = "trace_id"
	requestIDKey ctxKey = "request_id"
)

func TestLoggerInitialization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		opts          logger.Options
		expectedError bool
	}{
		{
			name: "Default options",
			opts: logger.DefaultOptions(),
		},
		{
			name: "Custom options",
			opts: logger.Options{
				Level:       logger.Debug,
				Encoding:    "json",
				OutputPath:  []string{"stdout"},
				ErrorPath:   []string{"stderr"},
				Development: true,
				LogDir:      "test_logs",
			},
		},
		{
			name: "Invalid log level",
			opts: logger.Options{
				Level:      "invalid",
				Encoding:   "json",
				OutputPath: []string{"stdout"},
				ErrorPath:  []string{"stderr"},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			log, err := logger.New(tt.opts)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, log)
			defer log.Close()
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	t.Parallel()

	tempFile := filepath.Join(os.TempDir(), "test_log.json")
	defer os.Remove(tempFile)

	opts := logger.Options{
		Level:      logger.Debug,
		Encoding:   "json",
		OutputPath: []string{tempFile},
		ErrorPath:  []string{tempFile},
	}

	log, err := logger.New(opts)
	require.NoError(t, err)
	defer log.Close()

	// Test different log levels
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	// Read the log file and verify content
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	logContent := string(content)

	// Verify all log levels are present
	assert.Contains(t, logContent, "debug message")
	assert.Contains(t, logContent, "info message")
	assert.Contains(t, logContent, "warn message")
	assert.Contains(t, logContent, "error message")
}

func TestGlobalLogger(t *testing.T) {
	t.Parallel()

	// Get global logger
	log1 := logger.GetLogger()
	assert.NotNil(t, log1)

	// Get another instance
	log2 := logger.GetLogger()

	// Verify that both instances are the same (singleton)
	assert.Equal(t, log1, log2)
}
