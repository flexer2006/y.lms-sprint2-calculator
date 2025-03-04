// Package storage предоставляет функции для управления выражениями в хранилище.
package storage

import (
	"fmt"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

// SaveExpression сохраняет выражение в памяти.
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
		zap.String(common.FieldID, expr.ID),
		zap.String(common.FieldExpression, expr.Expression),
		zap.String(common.FieldStatus, string(expr.Status)))
	return nil
}

// GetExpression извлекает выражение из хранилища по идентификатору.
func (s *Storage) GetExpression(id string) (*models.Expression, error) {
	if value, ok := s.expressions.Load(id); ok {
		s.logger.Debug("Expression retrieved",
			zap.String(common.FieldID, id))
		return value.(*models.Expression), nil
	}
	s.logger.Warn("Expression not found",
		zap.String(common.FieldID, id))
	return nil, fmt.Errorf("expression not found")
}

// UpdateExpressionStatus обновляет статус выражения в хранилище.
func (s *Storage) UpdateExpressionStatus(id string, status models.ExpressionStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)
		oldStatus := expr.Status

		if oldStatus == status {
			return nil
		}

		if !isValidStatusTransition(oldStatus, status) {
			s.logger.Error(common.LogInvalidStatusTransition,
				zap.String(common.FieldID, id),
				zap.String(common.FieldOldStatus, string(oldStatus)),
				zap.String(common.FieldNewStatus, string(status)))
			return fmt.Errorf("invalid status transition from %s to %s", oldStatus, status)
		}

		updated := *expr
		updated.Status = status
		updated.UpdatedAt = time.Now().Add(time.Millisecond)

		s.expressions.Store(id, &updated)
		s.logger.Info(common.LogExpressionStatusUpdated,
			zap.String(common.FieldID, id),
			zap.String(common.FieldOldStatus, string(oldStatus)),
			zap.String(common.FieldNewStatus, string(status)))
		return nil
	}
	s.logger.Error(common.LogFailedUpdateStatusNotFound,
		zap.String(common.FieldID, id))
	return fmt.Errorf("expression not found")
}

// UpdateExpressionResult обновляет результат выражения в хранилище.
func (s *Storage) UpdateExpressionResult(id string, result float64) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)

		updated := *expr
		updated.Result = &result
		updated.Status = models.StatusComplete
		updated.UpdatedAt = time.Now()

		s.expressions.Store(id, &updated)
		return nil
	}
	return fmt.Errorf("expression not found")
}

// UpdateExpressionError обновляет ошибку выражения в хранилище.
func (s *Storage) UpdateExpressionError(id string, err string) error {
	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)

		updated := *expr
		updated.Error = err
		updated.Status = models.StatusError
		updated.UpdatedAt = time.Now()

		s.expressions.Store(id, &updated)
		return nil
	}
	return fmt.Errorf("expression not found")
}

// ListExpressions перечисляет все выражения, находящиеся в хранилище.
func (s *Storage) ListExpressions() []*models.Expression {
	var expressions []*models.Expression
	s.expressions.Range(func(key, value interface{}) bool {
		expressions = append(expressions, value.(*models.Expression))
		return true
	})
	s.logger.Debug(common.LogListedAllExpressions,
		zap.Int(common.FieldCount, len(expressions)))
	return expressions
}

// isValidStatusTransition Проверяет, действителен ли переход состояния.
func isValidStatusTransition(from, to models.ExpressionStatus) bool {
	switch from {
	case models.StatusPending:
		return to == models.StatusProgress || to == models.StatusError
	case models.StatusProgress:
		return to == models.StatusComplete || to == models.StatusError
	case models.StatusComplete, models.StatusError:
		return false
	default:
		return true
	}
}
