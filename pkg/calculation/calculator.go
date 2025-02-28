// Package calculation provides functions to evaluate mathematical expressions.
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

	parser := &Parser{tokens: tokens, pos: 0}
	return parser.parse()
}
