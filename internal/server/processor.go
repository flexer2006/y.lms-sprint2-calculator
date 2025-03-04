// Package server предоставляет функциональные возможности для обработки математических выражений.
package server

import (
	"fmt"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// processExpression обрабатывает заданное математическое выражение, составляя задачи.
func (s *Server) processExpression(expr *models.Expression) error {
	tokens, err := s.parseExpression(expr.Expression)
	if err != nil {
		s.logger.Error("Failed to parse expression",
			zap.String("expression", expr.Expression),
			zap.Error(err))
		if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
			s.logger.Error("Failed to update expression error status", zap.Error(updateErr))
		}
		return err
	}

	if err := s.storage.UpdateExpressionStatus(expr.ID, models.StatusProgress); err != nil {
		s.logger.Error("Failed to update expression status", zap.Error(err))
		return err
	}

	tasks, err := s.createTasks(expr.ID, tokens)
	if err != nil {
		s.logger.Error("Failed to create tasks", zap.Error(err))
		if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
			s.logger.Error("Failed to update expression error status", zap.Error(updateErr))
		}
		return err
	}

	for _, task := range tasks {
		if err := s.storage.SaveTask(task); err != nil {
			s.logger.Error("Failed to save task", zap.Error(err))
			return err
		}
	}
	return nil
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
			return nil, fmt.Errorf("invalid expression: too few tokens")
		}
		return nil, fmt.Errorf("invalid expression: trailing operator")
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

// createTasks создает вычислительные задачи из лексем выражения.
func (s *Server) createTasks(exprID string, tokens []string) ([]*models.Task, error) {
	rpnTokens, err := s.toRPN(tokens)
	if err != nil {
		return nil, err
	}

	var tasks []*models.Task
	var stack []interface{}

	for _, token := range rpnTokens {
		if isOperator(token) {
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid RPN expression: too few operands")
			}

			op2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			op1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			task := &models.Task{
				ID:           uuid.New().String(),
				ExpressionID: exprID,
				Operation:    token,
			}

			_ = s.getOperationTime(token)

			switch v := op1.(type) {
			case float64:
				task.Arg1 = v
			case string:
				task.DependsOnTaskIDs = append(task.DependsOnTaskIDs, v)
			}
			switch v := op2.(type) {
			case float64:
				task.Arg2 = v
			case string:
				task.DependsOnTaskIDs = append(task.DependsOnTaskIDs, v)
			}

			tasks = append(tasks, task)
			stack = append(stack, task.ID)
		} else {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid RPN expression: too many operands")
	}

	return tasks, nil
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

// toRPN преобразует маркеры инфиксной нотации в обратную польскую нотацию
func (s *Server) toRPN(tokens []string) ([]string, error) {
	var stack []string
	var output []string

	for _, token := range tokens {
		if isOperator(token) {

			for len(stack) > 0 && isOperator(stack[len(stack)-1]) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {

			if _, err := strconv.ParseFloat(token, 64); err != nil {
				return nil, fmt.Errorf("invalid number: %s", token)
			}
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}
