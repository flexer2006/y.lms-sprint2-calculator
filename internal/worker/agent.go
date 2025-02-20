package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"go.uber.org/zap"
)

// Agent представляет собой агента-вычислителя
type Agent struct {
	config     *Config
	logger     *logger.Logger
	httpClient *http.Client
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// New создает нового агента
func New(cfg *Config, log *logger.Logger) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	return &Agent{
		config: cfg,
		logger: log,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start запускает агента
func (a *Agent) Start() error {
	a.logger.Info("Starting agent",
		zap.Int("computing_power", a.config.ComputingPower),
		zap.String("orchestrator_url", a.config.OrchestratorURL))

	for i := 0; i < a.config.ComputingPower; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}

	return nil
}

// Stop останавливает агента
func (a *Agent) Stop() {
	a.cancel()
	a.wg.Wait()
	a.logger.Info("Agent stopped")
}

// worker представляет собой горутину-вычислителя
func (a *Agent) worker(id int) {
	defer a.wg.Done()

	a.logger.Info("Starting worker", zap.Int("worker_id", id))

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Info("Worker stopped", zap.Int("worker_id", id))
			return
		default:
			if err := a.processTask(id); err != nil {
				a.logger.Error("Failed to process task",
					zap.Int("worker_id", id),
					zap.Error(err))
				time.Sleep(time.Second) // Небольшая задержка перед следующей попыткой
			}
		}
	}
}

// processTask обрабатывает одну задачу
func (a *Agent) processTask(workerID int) error {
	// Получаем задачу от оркестратора
	task, err := a.getTask()
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Если задач нет, ждем немного
	if task == nil {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	a.logger.Info("Processing task",
		zap.Int("worker_id", workerID),
		zap.String("task_id", task.ID),
		zap.String("operation", task.Operation))

	// Имитируем время выполнения операции
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	// Вычисляем результат
	result := a.Calculate(task)

	// Отправляем результат
	if err := a.sendResult(task.ID, result); err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	return nil
}

// getTask получает задачу от оркестратора
func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.httpClient.Get(fmt.Sprintf("%s/internal/task", a.config.OrchestratorURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Нет доступных задач
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return &taskResp.Task, nil
}

// sendResult отправляет результат вычисления оркестратору
func (a *Agent) sendResult(taskID string, result float64) error {
	taskResult := models.TaskResult{
		ID:     taskID,
		Result: result,
	}

	body, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Post(
		fmt.Sprintf("%s/internal/task", a.config.OrchestratorURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Calculate выполняет вычисление
func (a *Agent) Calculate(task *models.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			a.logger.Error("Division by zero",
				zap.String("task_id", task.ID))
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		a.logger.Error("Unknown operation",
			zap.String("task_id", task.ID),
			zap.String("operation", task.Operation))
		return 0
	}
}
