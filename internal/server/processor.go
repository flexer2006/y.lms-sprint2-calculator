// Package server provides functionalities for processing mathematical expressions.
package server

import (
	"fmt"
	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// processExpression processes a given mathematical expression by creating tasks.
func (s *Server) processExpression(expr *models.Expression) {
	s.logger.Info(common.LogProcessingExpression, zap.String("id", expr.ID), zap.String("expression", expr.Expression))

	if err := s.storage.UpdateExpressionStatus(expr.ID, models.StatusProgress); err != nil {
		s.logger.Error("Failed to update expression status",
			zap.String("id", expr.ID),
			zap.Error(err))
		return
	}

	tokens, err := s.parseExpression(expr.Expression)
	if err != nil {
		s.logger.Error(common.LogFailedParseExpression,
			zap.String("id", expr.ID),
			zap.String(common.FieldExpression, expr.Expression),
			zap.Error(err))

		if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
			s.logger.Error(common.ErrFailedUpdateExpr,
				zap.String("id", expr.ID),
				zap.Error(updateErr))
		}
		return
	}

	tasks := s.createTasks(expr.ID, tokens)
	if len(tasks) == 0 {
		s.logger.Error(common.LogNoValidTasksCreated,
			zap.String("id", expr.ID),
			zap.String(common.FieldExpression, expr.Expression))
		if updateErr := s.storage.UpdateExpressionError(expr.ID, "Failed to create valid tasks"); updateErr != nil {
			s.logger.Error(common.ErrFailedUpdateExpr,
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
				s.logger.Error(common.ErrFailedUpdateExpr,
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

// parseExpression parses a mathematical expression into tokens.
func (s *Server) parseExpression(expression string) ([]string, error) {
	if len(expression) == 0 {
		return nil, fmt.Errorf("invalid request body")
	}

	var (
		tokens     []string
		parenStack int
	)

	for i := 0; i < len(expression); i++ {
		c := expression[i]
		if c == ' ' {
			continue
		}

		if c == '(' {
			tokens = append(tokens, "(")
			parenStack++
			continue
		}
		if c == ')' {
			if parenStack == 0 {
				return nil, fmt.Errorf("invalid expression: unmatched parentheses")
			}
			tokens = append(tokens, ")")
			parenStack--
			continue
		}
		if isDigit(c) || c == '.' {
			j := i
			for j < len(expression) && (isDigit(expression[j]) || expression[j] == '.') {
				j++
			}
			tokens = append(tokens, expression[i:j])
			i = j - 1
			continue
		}
		if isOperator(string(c)) {
			if i > 0 && isOperator(string(expression[i-1])) && !(expression[i-1] == '(' && c == '-') {
				if c == '-' && expression[i-1] == '-' {
					return nil, fmt.Errorf("invalid expression: invalid structure")
				}
			}
			if c == '-' && (i == 0 || isOperator(string(expression[i-1])) || expression[i-1] == '(') {
				tokens = append(tokens, "-1", "*")
				continue
			}
			tokens = append(tokens, string(c))
			continue
		}
		return nil, fmt.Errorf("invalid expression: unexpected character '%c'", c)
	}

	if parenStack != 0 {
		return nil, fmt.Errorf("invalid expression: unmatched parentheses")
	}

	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i] == "(" && tokens[i+1] == ")" {
			return nil, fmt.Errorf("invalid expression: empty expression")
		}
	}

	for i := 0; i < len(tokens)-1; i++ {
		if isOperator(tokens[i]) && tokens[i+1] == ")" {
			return nil, fmt.Errorf("invalid expression: invalid structure")
		}
		if tokens[i] == "(" && isOperator(tokens[i+1]) {
			return nil, fmt.Errorf("invalid expression: invalid structure")
		}
	}

	operands, operators := 0, 0
	for _, token := range tokens {
		if isOperator(token) {
			operators++
			continue
		}
		if token != "(" && token != ")" {
			operands++
		}
	}

	if operators == 0 {
		return nil, fmt.Errorf("invalid expression: too few tokens")
	}

	if len(tokens) == 1 && isOperator(tokens[0]) {
		return nil, fmt.Errorf("invalid expression: too few tokens")
	}

	if len(tokens) > 1 && isOperator(tokens[len(tokens)-1]) {
		if len(tokens) == 2 {
			return nil, fmt.Errorf("invalid expression: too few tokens") // e.g., "2*"
		}
		return nil, fmt.Errorf("invalid expression: trailing operator") // e.g., "1+2+"
	}

	if operators > 0 && operands <= 1 {
		return nil, fmt.Errorf("invalid expression: too few tokens")
	}

	if operands <= operators && operators > 0 {
		return nil, fmt.Errorf("invalid expression: invalid structure")
	}

	for _, token := range tokens {
		if token != "(" && token != ")" && !isOperator(token) && strings.Count(token, ".") > 1 {
			return nil, fmt.Errorf("invalid expression: invalid number format")
		}
	}

	return tokens, nil
}

// createTasks creates computational tasks from tokens of an expression.
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

// getOperationTime returns the time required for a specific operation.
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

// isOperator checks if a token is a valid operator.
func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

// isDigit checks if a byte is a digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
