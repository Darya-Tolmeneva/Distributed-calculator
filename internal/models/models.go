package models

import (
	"errors"
	"time"
)

type Request struct {
	Expression string `json:"expression"`
}
type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type Expression struct {
	ID     int     `json:"id"`
	Expr   string  `json:"expression"`
	Status string  `json:"status"` // "pending", "in_progress", "completed", "failed"
	Result float64 `json:"result"`
}
type Task struct {
	ID            int           `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
	Result        float64       `json:"result"`
}

type Result struct {
	Result float64 `json:"result"`
	ID     int     `json:"id"`
}

var (
	ErrInvalidSymbols         = errors.New("invalid symbols")
	ErrInvalidParenthesis     = errors.New("invalid parenthesis")
	ErrInvalidOperations      = errors.New("invalid operations")
	ErrNotEnoughOperands      = errors.New("there are not enough operands in the expression")
	ErrDivisionByZero         = errors.New("division by zero")
	ErrExpressionNotEvaluated = errors.New("the expression is not evaluated")
)
