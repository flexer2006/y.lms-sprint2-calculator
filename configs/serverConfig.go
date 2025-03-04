// Package configs предоставляет конфигурационные структуры и функции для сервера.
package configs

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfig содержит всю конфигурацию сервера.
type ServerConfig struct {
	Port              string // Port на котором будет прослушиваться сервер.
	TimeAdditionMS    int64  // Время в миллисекундах для операций сложения.
	TimeSubtractionMS int64  // Время в миллисекундах для операций вычитания.
	TimeMultiplyMS    int64  // Время в миллисекундах для операций умножения.
	TimeDivisionMS    int64  // Время в миллисекундах для операций деления.
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

// getEnvString извлекает строковое значение из окружения или возвращает значение по умолчанию.
func getEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvInt64 извлекает значение int64 из среды или возвращает значение по умолчанию.
func getEnvInt64(key string, defaultValue int64) (int64, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
