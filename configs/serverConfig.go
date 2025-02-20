package configs

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfig holds all configuration for the server
type ServerConfig struct {
	Port              string
	TimeAdditionMS    int64
	TimeSubtractionMS int64
	TimeMultiplyMS    int64
	TimeDivisionMS    int64
}

// NewServerConfig creates a new ServerConfig instance
func NewServerConfig() (*ServerConfig, error) {
	timeAdd, err := getEnvInt64("TIME_ADDITION_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_ADDITION_MS: %w", err)
	}

	timeSub, err := getEnvInt64("TIME_SUBTRACTION_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_SUBTRACTION_MS: %w", err)
	}

	timeMul, err := getEnvInt64("TIME_MULTIPLICATIONS_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_MULTIPLICATIONS_MS: %w", err)
	}

	timeDiv, err := getEnvInt64("TIME_DIVISIONS_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_DIVISIONS_MS: %w", err)
	}

	return &ServerConfig{
		Port:              getEnvString("PORT", "8080"),
		TimeAdditionMS:    timeAdd,
		TimeSubtractionMS: timeSub,
		TimeMultiplyMS:    timeMul,
		TimeDivisionMS:    timeDiv,
	}, nil
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) (int64, error) {
	if value, exists := os.LookupEnv(key); exists {
		return strconv.ParseInt(value, 10, 64)
	}
	return defaultValue, nil
}
