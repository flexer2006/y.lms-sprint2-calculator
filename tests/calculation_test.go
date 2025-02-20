package test

import (
	"testing"

	"github.com/flexer2006/y.lms-sprint2-calculator/pkg/calculation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expr     string
		expected float64
		wantErr  bool
	}{

		{
			name:     "simple addition",
			expr:     "2 + 2",
			expected: 4,
		},
		{
			name:     "simple subtraction",
			expr:     "5 - 3",
			expected: 2,
		},
		{
			name:     "simple multiplication",
			expr:     "4 * 3",
			expected: 12,
		},
		{
			name:     "simple division",
			expr:     "10 / 2",
			expected: 5,
		},

		// Сложные выражения с приоритетами
		{
			name:     "operator precedence multiplication over addition",
			expr:     "2 + 3 * 4",
			expected: 14,
		},
		{
			name:     "operator precedence division over subtraction",
			expr:     "10 - 6 / 2",
			expected: 7,
		},
		{
			name:     "multiple operators with precedence",
			expr:     "2 + 3 * 4 - 5 / 2.5",
			expected: 12,
		},

		// Работа со скобками
		{
			name:     "simple parentheses",
			expr:     "(2 + 3) * 4",
			expected: 20,
		},
		{
			name:     "nested parentheses",
			expr:     "((2 + 3) * 2) + 1",
			expected: 11,
		},
		{
			name:     "complex nested parentheses",
			expr:     "(1 + (2 * (3 + 4) - 5)) * 2",
			expected: 20,
		},
		{
			name:     "multiple groups of parentheses",
			expr:     "(2 + 3) * (4 + 5)",
			expected: 45,
		},

		// Работа с отрицательными числами
		{
			name:     "negative number addition",
			expr:     "-2 + 3",
			expected: 1,
		},
		{
			name:     "adding negative numbers",
			expr:     "-2 + -3",
			expected: -5,
		},
		{
			name:     "subtracting negative number",
			expr:     "5 - (-3)",
			expected: 8,
		},
		{
			name:     "multiplying negative numbers",
			expr:     "-2 * -3",
			expected: 6,
		},
		{
			name:     "dividing negative numbers",
			expr:     "-6 / -2",
			expected: 3,
		},
		{
			name:     "complex expression with negatives",
			expr:     "-2 * (-3 + 4) - (-5 / -2)",
			expected: -4.5,
		},

		// Работа с десятичными числами
		{
			name:     "decimal addition",
			expr:     "2.5 + 3.7",
			expected: 6.2,
		},
		{
			name:     "decimal multiplication",
			expr:     "2.5 * 2.5",
			expected: 6.25,
		},
		{
			name:     "complex decimal expression",
			expr:     "2.5 * 3.2 + 4.8 / 2.4",
			expected: 10,
		},

		// Пробелы и форматирование
		{
			name:     "expression with spaces",
			expr:     "  2  +  2  ",
			expected: 4,
		},
		{
			name:     "expression with multiple spaces",
			expr:     "2   +   3   *   4",
			expected: 14,
		},

		// Крайние случаи и ошибки
		{
			name:    "empty expression",
			expr:    "",
			wantErr: true,
		},
		{
			name:    "only spaces",
			expr:    "     ",
			wantErr: true,
		},
		{
			name:    "unclosed parenthesis",
			expr:    "(2 + 3",
			wantErr: true,
		},
		{
			name:    "unopened parenthesis",
			expr:    "2 + 3)",
			wantErr: true,
		},
		{
			name:    "multiple decimal points",
			expr:    "2.2.2 + 1",
			wantErr: true,
		},
		{
			name:    "invalid character",
			expr:    "2 + a",
			wantErr: true,
		},
		{
			name:    "consecutive operators",
			expr:    "2 ++ 2",
			wantErr: true,
		},
		{
			name:    "division by zero",
			expr:    "2 / 0",
			wantErr: true,
		},
		{
			name:    "missing operand",
			expr:    "2 +",
			wantErr: true,
		},
		{
			name:    "missing operator",
			expr:    "2 2",
			wantErr: true,
		},

		// Очень длинные выражения
		{
			name:     "long expression",
			expr:     "1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10",
			expected: 55,
		},
		{
			name:     "long expression with mixed operators",
			expr:     "1 + 2 * 3 - 4 / 2 + 5 * (6 + 7) - 8 + 9",
			expected: 71,
		},

		// Унарный минус в разных контекстах
		{
			name:     "unary minus at start",
			expr:     "-2 * 3",
			expected: -6,
		},
		{
			name:     "unary minus after operator",
			expr:     "2 * -3",
			expected: -6,
		},
		{
			name:     "unary minus in parentheses",
			expr:     "2 * (-3)",
			expected: -6,
		},
		{
			name:     "multiple unary minuses",
			expr:     "2 * -(-3)",
			expected: 6,
		},
		{
			name:     "unary minus with decimals",
			expr:     "-2.5 * -3.2",
			expected: 8,
		},

		// Очень маленькие и большие числа
		{
			name:     "very small numbers",
			expr:     "0.0000001 + 0.0000002",
			expected: 0.0000003,
		},
		{
			name:     "very large numbers",
			expr:     "1000000 * 1000000",
			expected: 1000000000000,
		},

		// Комбинированные сложные случаи
		{
			name:     "complex mixed expression",
			expr:     "((-2.5 + 3.7) * (-2 + 4.2)) / (2 * -0.5)",
			expected: -2.64,
		},
		{
			name:     "deeply nested expression",
			expr:     "((((1 + 2) * 3) - 4) / 5) * (-2)",
			expected: -2,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculation.EvaluateExpression(tt.expr)

			if tt.wantErr {
				assert.Error(t, err, "Expected error for expression: %s", tt.expr)
				return
			}

			require.NoError(t, err, "Unexpected error for expression: %s", tt.expr)
			assert.InDelta(t, tt.expected, result, 1e-10, "Unexpected result for expression: %s", tt.expr)
		})
	}
}
