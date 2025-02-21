package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Server) processExpression(expr *models.Expression) {
	s.logger.Info("Processing expression", zap.String("id", expr.ID), zap.String("expression", expr.Expression))

	if err := s.storage.UpdateExpressionStatus(expr.ID, models.StatusProgress); err != nil {
		s.logger.Error("Failed to update expression status",
			zap.String("id", expr.ID),
			zap.Error(err))
		return
	}

	tokens, err := s.parseExpression(expr.Expression)
	if err != nil {
		s.logger.Error("Failed to parse expression",
			zap.String("id", expr.ID),
			zap.String("expression", expr.Expression),
			zap.Error(err))

		if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
			s.logger.Error("Failed to update expression error",
				zap.String("id", expr.ID),
				zap.Error(updateErr))
		}
		return
	}

	tasks := s.createTasks(expr.ID, tokens)
	if len(tasks) == 0 {
		s.logger.Error("No valid tasks created",
			zap.String("id", expr.ID),
			zap.String("expression", expr.Expression))
		if updateErr := s.storage.UpdateExpressionError(expr.ID, "Failed to create valid tasks"); updateErr != nil {
			s.logger.Error("Failed to update expression error",
				zap.String("id", expr.ID),
				zap.Error(updateErr))
		}
		return
	}

	for _, task := range tasks {
		if err := s.storage.SaveTask(task); err != nil {
			s.logger.Error("Failed to save task",
				zap.String("expressionID", expr.ID),
				zap.String("taskID", task.ID),
				zap.Error(err))
			if updateErr := s.storage.UpdateExpressionError(expr.ID, "Failed to create tasks"); updateErr != nil {
				s.logger.Error("Failed to update expression error",
					zap.String("id", expr.ID),
					zap.Error(updateErr))
			}
			return
		}
	}

	s.logger.Info("Tasks created successfully",
		zap.String("id", expr.ID),
		zap.Int("taskCount", len(tasks)))
}

func (s *Server) parseExpression(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	var tokens []string
	var currentNumber strings.Builder

	for i := 0; i < len(expression); i++ {
		char := expression[i]

		if char == '-' && (i == 0 || isOperator(string(expression[i-1]))) {
			currentNumber.WriteRune('-')
			continue
		}

		if isDigit(char) || char == '.' {
			currentNumber.WriteByte(char)

			if i == len(expression)-1 {
				tokens = append(tokens, currentNumber.String())
			}
			continue
		}

		if currentNumber.Len() > 0 {
			tokens = append(tokens, currentNumber.String())
			currentNumber.Reset()
		}

		if isOperator(string(char)) {
			tokens = append(tokens, string(char))
		}
	}

	if len(tokens) < 3 {
		return nil, fmt.Errorf("invalid expression: too few tokens")
	}

	return tokens, nil
}

func (s *Server) createTasks(exprID string, tokens []string) []*models.Task {
	var tasks []*models.Task

	for i := 0; i < len(tokens); i++ {
		if isOperator(tokens[i]) {
			operationTime := s.getOperationTime(tokens[i])

			task := &models.Task{
				ID:            uuid.New().String(),
				ExpressionID:  exprID,
				Operation:     tokens[i],
				OperationTime: operationTime,
			}

			if i > 0 && i < len(tokens)-1 {
				arg1, err1 := strconv.ParseFloat(tokens[i-1], 64)
				arg2, err2 := strconv.ParseFloat(tokens[i+1], 64)
				if err1 == nil && err2 == nil {
					task.Arg1 = arg1
					task.Arg2 = arg2
					tasks = append(tasks, task)
				}
			}
		}
	}

	return tasks
}

func (s *Server) getOperationTime(op string) int64 {
	switch op {
	case "+":
		return s.config.TimeAdditionMS
	case "-":
		return s.config.TimeSubtractionMS
	case "*":
		return s.config.TimeMultiplyMS
	case "/":
		return s.config.TimeDivisionMS
	default:
		return 100
	}
}

func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}