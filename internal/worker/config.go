package worker

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию агента
type Config struct {
	ComputingPower  int    // Количество воркеров
	OrchestratorURL string // URL оркестратора
}

// NewConfig создает новую конфигурацию агента
func NewConfig() (*Config, error) {
	power, err := getComputingPower()
	if err != nil {
		return nil, fmt.Errorf("failed to get computing power: %w", err)
	}

	orchestratorURL := getEnvString("ORCHESTRATOR_URL", "http://localhost:8080")

	return &Config{
		ComputingPower:  power,
		OrchestratorURL: orchestratorURL,
	}, nil
}

// getComputingPower получает количество воркеров из переменной окружения
func getComputingPower() (int, error) {
	powerStr := getEnvString("COMPUTING_POWER", "1")
	power, err := strconv.Atoi(powerStr)
	if err != nil {
		return 0, fmt.Errorf("invalid COMPUTING_POWER value: %w", err)
	}
	if power < 1 {
		return 0, fmt.Errorf("COMPUTING_POWER must be greater than 0")
	}
	return power, nil
}

// getEnvString получает значение переменной окружения с значением по умолчанию
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
