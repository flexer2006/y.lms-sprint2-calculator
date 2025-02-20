package configs

import (
	"fmt"
	"os"
	"strconv"
)

// WorkerConfig содержит конфигурацию агента
type WorkerConfig struct {
	ComputingPower  int    // Количество воркеров
	OrchestratorURL string // URL оркестратора
}

// NewWorkerConfig создает новую конфигурацию агента
func NewWorkerConfig() (*WorkerConfig, error) {
	power, err := getWorkerComputingPower()
	if err != nil {
		return nil, fmt.Errorf("failed to get computing power: %w", err)
	}

	orchestratorURL := getWorkerEnvString("ORCHESTRATOR_URL", "http://localhost:8080")

	return &WorkerConfig{
		ComputingPower:  power,
		OrchestratorURL: orchestratorURL,
	}, nil
}

// getWorkerComputingPower получает количество воркеров из переменной окружения
func getWorkerComputingPower() (int, error) {
	powerStr := getWorkerEnvString("COMPUTING_POWER", "1")
	power, err := strconv.Atoi(powerStr)
	if err != nil {
		return 0, fmt.Errorf("invalid COMPUTING_POWER value: %w", err)
	}
	if power < 1 {
		return 0, fmt.Errorf("COMPUTING_POWER must be greater than 0")
	}
	return power, nil
}

// getWorkerEnvString получает значение переменной окружения с значением по умолчанию
func getWorkerEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
