package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"go.uber.org/zap"
)

type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Изменяем на слайс для гарантированного FIFO порядка
	mu          sync.Mutex
	logger      *zap.Logger
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
		logger:    logger,
	}
}

func (s *Storage) SaveExpression(expr *models.Expression) error {
	if expr.ID == "" {
		s.logger.Error("Failed to save expression: empty ID")
		return fmt.Errorf("expression ID cannot be empty")
	}

	now := time.Now()
	if expr.CreatedAt.IsZero() {
		expr.CreatedAt = now
	}
	expr.UpdatedAt = now.Add(time.Millisecond)

	s.expressions.Store(expr.ID, expr)
	s.logger.Info("Expression saved successfully",
		zap.String("id", expr.ID),
		zap.String("expression", expr.Expression),
		zap.String("status", string(expr.Status)))
	return nil
}

func (s *Storage) GetExpression(id string) (*models.Expression, error) {
	if value, ok := s.expressions.Load(id); ok {
		s.logger.Debug("Expression retrieved",
			zap.String("id", id))
		return value.(*models.Expression), nil
	}
	s.logger.Warn("Expression not found",
		zap.String("id", id))
	return nil, fmt.Errorf("expression not found")
}

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

func (s *Storage) UpdateExpressionStatus(id string, status models.ExpressionStatus) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)
		oldStatus := expr.Status

		if !isValidStatusTransition(oldStatus, status) {
			s.logger.Error("Invalid status transition",
				zap.String("id", id),
				zap.String("oldStatus", string(oldStatus)),
				zap.String("newStatus", string(status)))
			return fmt.Errorf("invalid status transition from %s to %s", oldStatus, status)
		}

		updated := *expr
		updated.Status = status
		updated.UpdatedAt = time.Now().Add(time.Millisecond)

		s.expressions.Store(id, &updated)
		s.logger.Info("Expression status updated",
			zap.String("id", id),
			zap.String("oldStatus", string(oldStatus)),
			zap.String("newStatus", string(status)))
		return nil
	}
	s.logger.Error("Failed to update expression status: expression not found",
		zap.String("id", id))
	return fmt.Errorf("expression not found")
}

func (s *Storage) UpdateExpressionResult(id string, result float64) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)

		// Создаем копию для обновления
		updated := *expr
		updated.Result = &result
		updated.Status = models.StatusComplete
		updated.UpdatedAt = time.Now()

		// Сохраняем обновленную версию
		s.expressions.Store(id, &updated)
		return nil
	}
	return fmt.Errorf("expression not found")
}

func (s *Storage) UpdateExpressionError(id string, err string) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)

		// Создаем копию для обновления
		updated := *expr
		updated.Error = err
		updated.Status = models.StatusError
		updated.UpdatedAt = time.Now()

		// Сохраняем обновленную версию
		s.expressions.Store(id, &updated)
		return nil
	}
	return fmt.Errorf("expression not found")
}

func (s *Storage) ListExpressions() []*models.Expression {
	var expressions []*models.Expression
	s.expressions.Range(func(key, value interface{}) bool {
		expressions = append(expressions, value.(*models.Expression))
		return true
	})
	s.logger.Debug("Listed all expressions",
		zap.Int("count", len(expressions)))
	return expressions
}

func isValidStatusTransition(from, to models.ExpressionStatus) bool {
	switch from {
	case models.StatusPending:
		return to == models.StatusProgress || to == models.StatusError
	case models.StatusProgress:
		return to == models.StatusComplete || to == models.StatusError
	case models.StatusComplete, models.StatusError:
		return false // Нельзя изменить статус после завершения или ошибки
	default:
		return true // Разрешаем переход для новых выражений
	}
}
