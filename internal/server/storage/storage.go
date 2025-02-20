package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
)

type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Изменяем на слайс для гарантированного FIFO порядка
	mu          sync.Mutex
}

func New() *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
	}
}

func (s *Storage) SaveExpression(expr *models.Expression) error {
	if expr.ID == "" {
		return fmt.Errorf("expression ID cannot be empty")
	}

	now := time.Now()
	if expr.CreatedAt.IsZero() {
		expr.CreatedAt = now
	}
	expr.UpdatedAt = now.Add(time.Millisecond) // Ensure UpdatedAt is after CreatedAt

	s.expressions.Store(expr.ID, expr)
	return nil
}

func (s *Storage) GetExpression(id string) (*models.Expression, error) {
	if value, ok := s.expressions.Load(id); ok {
		return value.(*models.Expression), nil
	}
	return nil, fmt.Errorf("expression not found")
}

func (s *Storage) ListExpressions() []models.Expression {
	var result []models.Expression
	s.expressions.Range(func(key, value interface{}) bool {
		result = append(result, *value.(*models.Expression))
		return true
	})
	return result
}

func (s *Storage) SaveTask(task *models.Task) error {
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	task.CreatedAt = now

	// Создаем копию задачи для хранения
	taskCopy := *task

	// Сохраняем в Map
	s.tasks.Store(task.ID, &taskCopy)

	// Добавляем в очередь
	s.taskQueue = append(s.taskQueue, taskCopy)

	return nil
}

func (s *Storage) GetTask(id string) (*models.Task, error) {
	if value, ok := s.tasks.Load(id); ok {
		return value.(*models.Task), nil
	}
	return nil, fmt.Errorf("task not found")
}

func (s *Storage) UpdateTaskResult(id string, result float64) error {
	if value, ok := s.tasks.Load(id); ok {
		task := value.(*models.Task)
		task.Result = &result
		s.tasks.Store(id, task)
		return nil
	}
	return fmt.Errorf("task not found")
}

func (s *Storage) GetNextTask() (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.taskQueue) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	// Получаем первую задачу из очереди
	task := s.taskQueue[0]

	// Удаляем задачу из очереди
	s.taskQueue = s.taskQueue[1:]

	// Возвращаем копию задачи
	return &task, nil
}

func (s *Storage) UpdateExpressionStatus(id string, status models.ExpressionStatus) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)
		oldStatus := expr.Status

		// Проверяем валидность перехода статуса
		if !isValidStatusTransition(oldStatus, status) {
			return fmt.Errorf("invalid status transition from %s to %s", oldStatus, status)
		}

		// Создаем копию для обновления
		updated := *expr
		updated.Status = status
		updated.UpdatedAt = time.Now().Add(time.Millisecond) // Гарантируем, что UpdatedAt будет позже CreatedAt

		// Сохраняем обновленную версию
		s.expressions.Store(id, &updated)
		return nil
	}
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
