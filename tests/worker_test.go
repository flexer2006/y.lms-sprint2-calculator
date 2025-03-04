package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/worker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgent_Calculate(t *testing.T) {
	t.Parallel()
	log, err := logger.New(logger.Options{
		Level:       logger.Debug,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: true,
	})
	require.NoError(t, err)

	agent := worker.New(&configs.WorkerConfig{ComputingPower: 1}, log)

	tests := []struct {
		name        string
		task        *models.Task
		expected    float64
		expectError bool
	}{
		{
			name: "Addition",
			task: &models.Task{
				ID:               "1",
				Operation:        "+",
				Arg1:             10,
				Arg2:             5,
				DependsOnTaskIDs: []string{}, // Добавлено для соответствия новой структуре, но не используется в Calculate
			},
			expected:    15,
			expectError: false,
		},
		{
			name: "Subtraction",
			task: &models.Task{
				ID:               "2",
				Operation:        "-",
				Arg1:             10,
				Arg2:             5,
				DependsOnTaskIDs: []string{},
			},
			expected:    5,
			expectError: false,
		},
		{
			name: "Multiplication",
			task: &models.Task{
				ID:               "3",
				Operation:        "*",
				Arg1:             10,
				Arg2:             5,
				DependsOnTaskIDs: []string{},
			},
			expected:    50,
			expectError: false,
		},
		{
			name: "Division",
			task: &models.Task{
				ID:               "4",
				Operation:        "/",
				Arg1:             10,
				Arg2:             5,
				DependsOnTaskIDs: []string{},
			},
			expected:    2,
			expectError: false,
		},
		{
			name: "Division by zero",
			task: &models.Task{
				ID:               "5",
				Operation:        "/",
				Arg1:             10,
				Arg2:             0,
				DependsOnTaskIDs: []string{},
			},
			expectError: true,
		},
		{
			name: "Unknown operation",
			task: &models.Task{
				ID:               "6",
				Operation:        "%",
				Arg1:             10,
				Arg2:             5,
				DependsOnTaskIDs: []string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				assert.Panics(t, func() { agent.Calculate(tt.task) }, "Expected panic but got none")
			} else {
				result := agent.Calculate(tt.task)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAgent_Integration(t *testing.T) {
	taskCh := make(chan models.Task, 1)
	resultCh := make(chan models.TaskResult, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			select {
			case task := <-taskCh:
				resp := models.TaskResponse{Task: task}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPost:
			var result models.TaskResult
			if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			select {
			case resultCh <- result:
				w.WriteHeader(http.StatusOK)
			case <-time.After(2 * time.Second):
				w.WriteHeader(http.StatusRequestTimeout)
			}
		}
	}))
	defer server.Close()

	log, err := logger.New(logger.Options{
		Level:       logger.Debug,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: true,
	})
	require.NoError(t, err)

	agent := worker.New(&configs.WorkerConfig{
		ComputingPower:  1,
		OrchestratorURL: server.URL,
	}, log)

	errCh := make(chan error, 1)
	go func() {
		errCh <- agent.Start()
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout starting agent")
	}

	defer func() {
		agent.Stop()

		select {
		case <-errCh:
		case <-time.After(5 * time.Second):
			t.Log("warning: agent stop timeout")
		}
	}()

	task := models.Task{
		ID:               "test-task",
		Operation:        "+",
		Arg1:             10,
		Arg2:             5,
		DependsOnTaskIDs: []string{}, // Добавлено для соответствия новой структуре
	}

	select {
	case taskCh <- task:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout sending task")
	}

	select {
	case result := <-resultCh:
		assert.Equal(t, task.ID, result.ID)
		assert.Equal(t, float64(15), result.Result)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for result")
	}
}

func TestAgent_Config(t *testing.T) {
	t.Parallel()

	originalComputingPower := os.Getenv("COMPUTING_POWER")
	originalOrchestratorURL := os.Getenv("ORCHESTRATOR_URL")

	defer func() {
		_ = os.Setenv("COMPUTING_POWER", originalComputingPower)
		_ = os.Setenv("ORCHESTRATOR_URL", originalOrchestratorURL)
	}()

	require.NoError(t, os.Setenv("COMPUTING_POWER", "3"))
	require.NoError(t, os.Setenv("ORCHESTRATOR_URL", "http://test:8080"))

	cfg, err := configs.NewWorkerConfig()
	require.NoError(t, err)

	assert.Equal(t, 3, cfg.ComputingPower)
	assert.Equal(t, "http://test:8080", cfg.OrchestratorURL)
}

func TestAgent_InvalidConfig(t *testing.T) {
	t.Parallel()

	originalComputingPower := os.Getenv("COMPUTING_POWER")
	originalOrchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	defer func() {
		if err := os.Setenv("COMPUTING_POWER", originalComputingPower); err != nil {
			t.Errorf("Failed to restore COMPUTING_POWER: %v", err)
		}
		if err := os.Setenv("ORCHESTRATOR_URL", originalOrchestratorURL); err != nil {
			t.Errorf("Failed to restore ORCHESTRATOR_URL: %v", err)
		}
	}()

	if err := os.Setenv("COMPUTING_POWER", "invalid"); err != nil {
		t.Fatalf("Failed to set COMPUTING_POWER: %v", err)
	}
	if err := os.Setenv("ORCHESTRATOR_URL", "http://test:8080"); err != nil {
		t.Fatalf("Failed to set ORCHESTRATOR_URL: %v", err)
	}

	_, err := configs.NewWorkerConfig()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid COMPUTING_POWER value")
	}

	if err := os.Setenv("COMPUTING_POWER", "0"); err != nil {
		t.Fatalf("Failed to set COMPUTING_POWER: %v", err)
	}
	_, err = configs.NewWorkerConfig()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "COMPUTING_POWER must be greater than 0")
	}
}
