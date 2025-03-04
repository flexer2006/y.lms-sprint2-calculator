// Package configs предоставляет конфигурационные структуры и функции для рабочего агента.
package configs

import (
	"fmt"
	"os"
	"strconv"
)

// WorkerConfig содержит конфигурацию рабочего агента.
type WorkerConfig struct {
	ComputingPower    int    // Количество рабочих.
	OrchestratorURL   string // URL-адрес оркестратора.
	AdditionTimeMS    int64  // Время в миллисекундах для операций сложения.
	SubtractionTimeMS int64  // Время в миллисекундах для операций вычитания.
	MultiplyTimeMS    int64  // Время в миллисекундах для операций умножения.
	DivisionTimeMS    int64  // Время в миллисекундах для операций деления.
}

// NewWorkerConfig создает новую конфигурацию рабочего агента.
func NewWorkerConfig() (*WorkerConfig, error) {
	power, err := getWorkerComputingPower()
	if err != nil {
		return nil, fmt.Errorf("failed to get computing power: %w", err)
	}

	timeAdd, err := getWorkerEnvInt64("TIME_ADDITION_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_ADDITION_MS: %w", err)
	}

	timeSub, err := getWorkerEnvInt64("TIME_SUBTRACTION_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_SUBTRACTION_MS: %w", err)
	}

	timeMul, err := getWorkerEnvInt64("TIME_MULTIPLICATIONS_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_MULTIPLICATIONS_MS: %w", err)
	}

	timeDiv, err := getWorkerEnvInt64("TIME_DIVISIONS_MS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid TIME_DIVISIONS_MS: %w", err)
	}

	return &WorkerConfig{
		ComputingPower:    power,
		OrchestratorURL:   getWorkerEnvString("ORCHESTRATOR_URL", "http://localhost:8080"),
		AdditionTimeMS:    timeAdd,
		SubtractionTimeMS: timeSub,
		MultiplyTimeMS:    timeMul,
		DivisionTimeMS:    timeDiv,
	}, nil
}

// getWorkerComputingPower извлекает количество рабочих из переменной окружения.
func getWorkerComputingPower() (int, error) {
	powerStr := getWorkerEnvString("COMPUTING_POWER", "1")

	power, err := strconv.Atoi(powerStr)
	if err != nil {
		return 0, fmt.Errorf("invalid COMPUTING_POWER value: %s", powerStr)
	}

	if power < 1 {
		return 0, fmt.Errorf("COMPUTING_POWER must be greater than 0")
	}

	return power, nil
}

// getWorkerEnvString извлекает значение переменной окружения со значением по умолчанию.
func getWorkerEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

// getWorkerEnvInt64 извлекает значение int64 из окружения или возвращает значение по умолчанию.
func getWorkerEnvInt64(key string, defaultValue int64) (int64, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
