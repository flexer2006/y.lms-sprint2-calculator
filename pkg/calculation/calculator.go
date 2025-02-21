package calculation

import (
	"errors"

	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger initializes the logger for the calculation package
func InitLogger(l *zap.Logger) {
	logger = l
}

// EvaluateExpression evaluates a mathematical expression and returns the result
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
