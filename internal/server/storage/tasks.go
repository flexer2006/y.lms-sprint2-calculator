package storage

import (
	"fmt"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

func (s *Storage) SaveTask(task *models.Task) error {
	if task.ID == "" {
		s.logger.Error("Failed to save task: empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	task.CreatedAt = now

	taskCopy := *task
	s.tasks.Store(task.ID, &taskCopy)
	s.taskQueue = append(s.taskQueue, taskCopy)

	s.logger.Info("Task saved successfully",
		zap.String("id", task.ID),
		zap.String("expressionID", task.ExpressionID),
		zap.String("operation", task.Operation))
	return nil
}

func (s *Storage) GetTask(id string) (*models.Task, error) {
	if value, ok := s.tasks.Load(id); ok {
		s.logger.Debug("Task retrieved",
			zap.String("id", id))
		return value.(*models.Task), nil
	}
	s.logger.Warn("Task not found",
		zap.String("id", id))
	return nil, fmt.Errorf("task not found")
}

func (s *Storage) UpdateTaskResult(id string, result float64) error {
	if value, ok := s.tasks.Load(id); ok {
		task := value.(*models.Task)
		task.Result = &result
		s.tasks.Store(id, task)
		s.logger.Info("Task result updated",
			zap.String("id", id),
			zap.Float64("result", result))
		return nil
	}
	s.logger.Error("Failed to update task result: task not found",
		zap.String("id", id))
	return fmt.Errorf("task not found")
}

func (s *Storage) GetNextTask() (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.taskQueue) == 0 {
		s.logger.Debug("No tasks available in queue")
		return nil, fmt.Errorf("no tasks available")
	}

	task := s.taskQueue[0]
	s.taskQueue = s.taskQueue[1:]

	s.logger.Info("Next task retrieved from queue",
		zap.String("id", task.ID),
		zap.String("expressionID", task.ExpressionID))
	return &task, nil
}
