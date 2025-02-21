package test

import (
	"testing"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			defer func() {
				if err := log.Close(); err != nil {
					t.Errorf("Failed to close logger: %v", err)
				}
			}()
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	t.Parallel()

	opts := logger.Options{
		Level:      logger.Debug,
		Encoding:   "json",
		OutputPath: []string{"stdout"},
		ErrorPath:  []string{"stderr"},
	}

	log, err := logger.New(opts)
	require.NoError(t, err)

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	require.NoError(t, log.Close())
}

func TestGlobalLogger(t *testing.T) {
	t.Parallel()

	log1 := logger.GetLogger()
	assert.NotNil(t, log1)

	log2 := logger.GetLogger()

	assert.Equal(t, log1, log2)
}
