// Package calculation предоставляет функции для оценки математических выражений.
package calculation

import (
	"errors"

	"go.uber.org/zap"
)

var logger *zap.Logger

// EvaluateExpression evaluates a mathematical expression and returns the result.
// It returns an error if the expression is empty or invalid.
func EvaluateExpression(expression string) (float64, error) {
	if expression == "" {
		return 0, errors.New("expression is empty")
	}

	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("invalid expression")
	}

	if logger != nil {
		logger.Debug("Tokens generated", zap.Strings("tokens", tokens))
	}

	parser := &Parser{tokens: tokens, pos: 0}
	result, err := parser.parse()
	if err != nil {
		if logger != nil {
			logger.Error("Parser failed", zap.Error(err), zap.String("expression", expression))
		}
		return 0, err
	}
	return result, nil
}
