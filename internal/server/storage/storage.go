// Package storage предоставляет функции хранения данных для выражений и задач.
package storage

import (
	"fmt"
	"sync"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

// Storage управляет хранением выражений и задач.
type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Slice to ensure FIFO order
	mu          sync.Mutex
	logger      *zap.Logger
}

// New создает новый экземпляр Storage с предоставленным logger.
func New(logger *zap.Logger) *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
		logger:    logger,
	}
}

// GetTasksByDependency извлекает задачи, зависящие от заданного идентификатора задачи.
func (s *Storage) GetTasksByDependency(taskID string) []*models.Task {
	var dependentTasks []*models.Task
	s.tasks.Range(func(_, value interface{}) bool {
		task := value.(*models.Task)
		for _, depID := range task.DependsOnTaskIDs {
			if depID == taskID {
				dependentTasks = append(dependentTasks, task)
				break
			}
		}
		return true
	})
	return dependentTasks
}

// GetTaskResult получает результат задачи по идентификатору.
func (s *Storage) GetTaskResult(taskID string) (float64, error) {
	if value, ok := s.tasks.Load(taskID); ok {
		task := value.(*models.Task)
		if task.Result == nil {
			return 0, fmt.Errorf("task result not set: %s", taskID)
		}
		return *task.Result, nil
	}
	return 0, fmt.Errorf("task not found") // Исправлено на константную строку вместо strings.ToLower
}

// GetTasksByExpressionID извлекает все задачи, связанные с идентификатором выражения.
func (s *Storage) GetTasksByExpressionID(expressionID string) []*models.Task {
	var tasks []*models.Task
	s.tasks.Range(func(_, value interface{}) bool {
		task := value.(*models.Task)
		if task.ExpressionID == expressionID {
			tasks = append(tasks, task)
		}
		return true
	})
	return tasks
}
