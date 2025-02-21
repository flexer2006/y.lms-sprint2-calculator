// Package models defines the data structures used in the server.
package models

import (
	"time"
)

// ExpressionStatus represents the status of an expression.
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

// Expression represents a mathematical expression and its status.
type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression,omitempty"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	Error      string           `json:"error,omitempty"`
}

// Task represents a computational task with two arguments and an operation.
type Task struct {
	ID            string    `json:"id"`
	Arg1          float64   `json:"arg1"`
	Arg2          float64   `json:"arg2"`
	Operation     string    `json:"operation"`
	OperationTime int64     `json:"operation_time"`
	ExpressionID  string    `json:"-"`
	CreatedAt     time.Time `json:"-"`
	Result        *float64  `json:"result,omitempty"`
}

// CalculateRequest represents a request to calculate an expression.
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse represents the response containing the ID of the calculation.
type CalculateResponse struct {
	ID string `json:"id"`
}

// TaskResult represents the result of a task calculation.
type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

// ExpressionResponse represents a response containing a single expression.
type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

// ExpressionsResponse represents a response containing multiple expressions.
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

// TaskResponse represents a response containing a single task.
type TaskResponse struct {
	Task Task `json:"task"`
}
