package models

import (
	"errors"
	"log"
	"os"
	"strconv"
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

var (
	TIME_ADDITION_MS        time.Duration
	TIME_SUBTRACTION_MS     time.Duration
	TIME_MULTIPLICATIONS_MS time.Duration
	TIME_DIVISIONS_MS       time.Duration
)

func init() {
	TIME_ADDITION_MS = parseDuration("TIME_ADDITION_MS", 1000*time.Millisecond)
	TIME_SUBTRACTION_MS = parseDuration("TIME_SUBTRACTION_MS", 1000*time.Millisecond)
	TIME_MULTIPLICATIONS_MS = parseDuration("TIME_MULTIPLICATIONS_MS", 1000*time.Millisecond)
	TIME_DIVISIONS_MS = parseDuration("TIME_DIVISIONS_MS", 1000*time.Millisecond)
}

func parseDuration(envVar string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", envVar, err)
	}
	return time.Duration(value) * time.Millisecond
}
