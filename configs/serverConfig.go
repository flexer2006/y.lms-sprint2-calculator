// Package configs provides configuration structures and functions for the server.
package configs

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfig holds all configuration for the server.
type ServerConfig struct {
	Port              string // Port on which the server will listen.
	TimeAdditionMS    int64  // Time in milliseconds for addition operations.
	TimeSubtractionMS int64  // Time in milliseconds for subtraction operations.
	TimeMultiplyMS    int64  // Time in milliseconds for multiplication operations.
	TimeDivisionMS    int64  // Time in milliseconds for division operations.
}

// NewServerConfig creates a new ServerConfig instance with values from environment variables or defaults.
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

	port := getEnvString("PORT", "8080")

	return &ServerConfig{
		Port:              port,
		TimeAdditionMS:    timeAdd,
		TimeSubtractionMS: timeSub,
		TimeMultiplyMS:    timeMul,
		TimeDivisionMS:    timeDiv,
	}, nil
}

// getEnvString retrieves a string value from the environment or returns a default value.
func getEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvInt64 retrieves an int64 value from the environment or returns a default value.
func getEnvInt64(key string, defaultValue int64) (int64, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
