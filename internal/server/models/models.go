// Package models определяет структуры данных, используемые в сервере.
package models

import (
	"time"
)

// ExpressionStatus представляет собой статус выражения.
type ExpressionStatus string

const (
	// StatusPending indicates the expression is pending processing.
	StatusPending ExpressionStatus = "PENDING"
	// StatusProgress indicates the expression is currently being processed.
	StatusProgress ExpressionStatus = "IN_PROGRESS"
	// StatusComplete indicates the expression has been processed successfully.
	StatusComplete ExpressionStatus = "COMPLETE"
	// StatusError indicates there was an error processing the expression.
	StatusError ExpressionStatus = "ERROR"
)

// Expression представляет собой математическое выражение и его состояние.
type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression,omitempty"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	Error      string           `json:"error,omitempty"`
}

// Task представляет собой вычислительную задачу с двумя аргументами и операцией.
type Task struct {
	ID               string
	ExpressionID     string
	Operation        string
	Arg1             float64
	Arg2             float64
	Result           *float64 // nil
	CreatedAt        time.Time
	DependsOnTaskIDs []string
}

// CalculateRequest представляет собой запрос на вычисление выражения.
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse представляет собой ответ, содержащий идентификатор вычисления.
type CalculateResponse struct {
	ID string `json:"id"`
}

// TaskResult представляет собой результат вычисления задачи.
type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

// ExpressionResponse представляет собой ответ, содержащий одно выражение.
type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

// ExpressionsResponse представляет собой ответ, содержащий несколько выражений.
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

// TaskResponse представляет собой ответ, содержащий одно задание.
type TaskResponse struct {
	Task Task `json:"task"`
}
