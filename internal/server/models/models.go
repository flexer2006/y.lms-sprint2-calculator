package models

import (
	"time"
)

type ExpressionStatus string

const (
	StatusPending  ExpressionStatus = "PENDING"
	StatusProgress ExpressionStatus = "IN_PROGRESS"
	StatusComplete ExpressionStatus = "COMPLETE"
	StatusError    ExpressionStatus = "ERROR"
)

type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression,omitempty"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	Error      string           `json:"error,omitempty"`
}

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

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID string `json:"id"`
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type TaskResponse struct {
	Task Task `json:"task"`
}
