package worker

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

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
