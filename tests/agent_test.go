package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgent_Calculate(t *testing.T) {
	log, err := logger.New(logger.Options{
		Level:       logger.Debug,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: true,
	})
	require.NoError(t, err)

	agent := worker.New(&worker.Config{ComputingPower: 1}, log)

	tests := []struct {
		name     string
		task     *models.Task
		expected float64
	}{
		{
			name: "Addition",
			task: &models.Task{
				ID:        "1",
				Operation: "+",
				Arg1:      10,
				Arg2:      5,
			},
			expected: 15,
		},
		{
			name: "Subtraction",
			task: &models.Task{
				ID:        "2",
				Operation: "-",
				Arg1:      10,
				Arg2:      5,
			},
			expected: 5,
		},
		{
			name: "Multiplication",
			task: &models.Task{
				ID:        "3",
				Operation: "*",
				Arg1:      10,
				Arg2:      5,
			},
			expected: 50,
		},
		{
			name: "Division",
			task: &models.Task{
				ID:        "4",
				Operation: "/",
				Arg1:      10,
				Arg2:      5,
			},
			expected: 2,
		},
		{
			name: "Division by zero",
			task: &models.Task{
				ID:        "5",
				Operation: "/",
				Arg1:      10,
				Arg2:      0,
			},
			expected: 0,
		},
		{
			name: "Unknown operation",
			task: &models.Task{
				ID:        "6",
				Operation: "%",
				Arg1:      10,
				Arg2:      5,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.Calculate(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAgent_Integration(t *testing.T) {
	// Создаем тестовый сервер
	taskCh := make(chan models.Task, 1)
	resultCh := make(chan models.TaskResult, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			select {
			case task := <-taskCh:
				resp := models.TaskResponse{Task: task}
				json.NewEncoder(w).Encode(resp)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPost:
			var result models.TaskResult
			json.NewDecoder(r.Body).Decode(&result)
			resultCh <- result
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// Создаем агента
	log, err := logger.New(logger.Options{
		Level:       logger.Debug,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: true,
	})
	require.NoError(t, err)

	agent := worker.New(&worker.Config{
		ComputingPower:  1,
		OrchestratorURL: server.URL,
	}, log)

	// Запускаем агента
	err = agent.Start()
	require.NoError(t, err)
	defer agent.Stop()

	// Отправляем задачу
	task := models.Task{
		ID:            "test-task",
		Operation:     "+",
		Arg1:          10,
		Arg2:          5,
		OperationTime: 100,
	}
	taskCh <- task

	// Ждем результат
	select {
	case result := <-resultCh:
		assert.Equal(t, task.ID, result.ID)
		assert.Equal(t, float64(15), result.Result)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for result")
	}
}

func TestAgent_Config(t *testing.T) {
	t.Setenv("COMPUTING_POWER", "3")
	t.Setenv("ORCHESTRATOR_URL", "http://test:8080")

	cfg, err := worker.NewConfig()
	require.NoError(t, err)

	assert.Equal(t, 3, cfg.ComputingPower)
	assert.Equal(t, "http://test:8080", cfg.OrchestratorURL)
}

func TestAgent_InvalidConfig(t *testing.T) {
	t.Setenv("COMPUTING_POWER", "invalid")

	_, err := worker.NewConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid COMPUTING_POWER value")

	t.Setenv("COMPUTING_POWER", "0")
	_, err = worker.NewConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMPUTING_POWER must be greater than 0")
}
