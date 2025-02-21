package storage

import (
	"fmt"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

func (s *Storage) SaveExpression(expr *models.Expression) error {
	if expr.ID == "" {
		s.logger.Error(common.LogFailedSaveEmptyID)
		return fmt.Errorf(common.ErrEmptyExpressionID)
	}

	now := time.Now()
	if expr.CreatedAt.IsZero() {
		expr.CreatedAt = now
	}
	expr.UpdatedAt = now.Add(time.Millisecond)

	s.expressions.Store(expr.ID, expr)
	s.logger.Info(common.LogExpressionSaved,
		zap.String(common.FieldID, expr.ID),
		zap.String(common.FieldExpression, expr.Expression),
		zap.String(common.FieldStatus, string(expr.Status)))
	return nil
}

func (s *Storage) GetExpression(id string) (*models.Expression, error) {
	if value, ok := s.expressions.Load(id); ok {
		s.logger.Debug(common.LogExpressionRetrieved,
			zap.String(common.FieldID, id))
		return value.(*models.Expression), nil
	}
	s.logger.Warn(common.LogExpressionNotFound,
		zap.String(common.FieldID, id))
	return nil, fmt.Errorf(common.ErrExpressionNotFound)
}

func (s *Storage) UpdateExpressionStatus(id string, status models.ExpressionStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value, ok := s.expressions.Load(id); ok {
		expr := value.(*models.Expression)
		oldStatus := expr.Status

		if oldStatus == status {
			return nil // Already in desired status
		}

		if !isValidStatusTransition(oldStatus, status) {
			s.logger.Error(common.LogInvalidStatusTransition,
				zap.String(common.FieldID, id),
				zap.String(common.FieldOldStatus, string(oldStatus)),
				zap.String(common.FieldNewStatus, string(status)))
			return fmt.Errorf(common.ErrInvalidStatusTransition, oldStatus, status)
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
	return fmt.Errorf(common.ErrExpressionNotFoundStorage)
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
	return fmt.Errorf(common.ErrExpressionNotFoundStorage)
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
	return fmt.Errorf(common.ErrExpressionNotFoundStorage)
}

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
