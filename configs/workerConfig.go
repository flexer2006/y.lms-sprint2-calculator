// Package configs provides configuration structures and functions for the worker agent.
package configs

import (
	"fmt"
	"os"
	"strconv"
)

// WorkerConfig contains the configuration for the worker agent.
type WorkerConfig struct {
	ComputingPower  int    // Number of workers.
	OrchestratorURL string // URL of the orchestrator.
}

// NewWorkerConfig creates a new worker agent configuration.
func NewWorkerConfig() (*WorkerConfig, error) {
	power, err := getWorkerComputingPower()
	if err != nil {
		return nil, fmt.Errorf("failed to get computing power: %w", err)
	}

	return &WorkerConfig{
		ComputingPower:  power,
		OrchestratorURL: getWorkerEnvString("ORCHESTRATOR_URL", "http://localhost:8080"),
	}, nil
}

// getWorkerComputingPower retrieves the number of workers from the environment variable.
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

// getWorkerEnvString retrieves an environment variable value with a default value.
func getWorkerEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
