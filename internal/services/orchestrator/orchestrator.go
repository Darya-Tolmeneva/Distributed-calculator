package orchestrator

import (
	"Distributed_calculator/internal/models"
	"Distributed_calculator/internal/services/checker"
	"Distributed_calculator/internal/services/parser"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	agentAddresses = []string{"http://localhost:50050", "http://localhost:50051", "http://localhost:50052", "http://localhost:50053", "http://localhost:50054"}
	currentAgent   int
	mu             sync.Mutex
	tasks          []*models.Task
	taskQueueMu    sync.Mutex
	taskId         int
	expressions    = make(map[int]models.Expression)
	expressionsMu  sync.Mutex
	previousId     int
)

// getNextAgent возвращает адрес следующего агента для распределения задач
func getNextAgent() string {
	mu.Lock()
	defer mu.Unlock()
	agent := agentAddresses[currentAgent]
	currentAgent = (currentAgent + 1) % len(agentAddresses)
	return agent
}

// evaluateTask отправляет задачу агенту и возвращает результат
func evaluateTask(task *models.Task) error {
	agentAddr := getNextAgent()
	log.Printf("Take agent %s", agentAddr)

	payload, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	resp, err := http.Post(agentAddr+"/compute", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send task to agent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent returned status code %d", resp.StatusCode)
	}
	return nil
}

func GetTaskResult(result models.Result) {
	taskQueueMu.Lock()
	defer taskQueueMu.Unlock()
	for _, task := range tasks {
		if task.ID == result.ID {
			task.Result = result.Result
		}
	}
}

// evaluateExpression вычисляет выражение
func evaluateExpression(expression models.Expression) (float64, error) {
	ch := checker.ExpressionChecker{}
	if err := checker.ValidateExpression(ch, expression.Expr); err != nil {
		return 0, fmt.Errorf("invalid expression: %v", err)
	}
	log.Println("Check expression")

	tasksPrefix := parser.ParseToPostfix(expression.Expr)
	log.Printf("Expression was parsed: %v", tasksPrefix)

	var stack []float64
	for _, token := range tasksPrefix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
			continue
		}
		if len(token) == 1 && checker.IsOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, models.ErrNotEnoughOperands
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if rune(token[0]) == '/' && b == 0 {
				return 0, models.ErrDivisionByZero
			}
			task := &models.Task{
				ID:        taskId,
				Arg1:      a,
				Arg2:      b,
				Operation: string(token[0]),
				Result:    float64(0),
			}
			tasks = append(tasks, task)
			log.Println(task.ID)
			err := evaluateTask(task)

			if err != nil {
				return 0, fmt.Errorf("failed to evaluate task: %v", err)
			}
			result := task.Result
			stack = append(stack, result)
			taskId++

		}
	}

	if len(stack) != 1 {
		return 0, models.ErrExpressionNotEvaluated
	}
	updateExpressionStatus(expression.ID, "done", stack[0])
	return stack[0], nil
}

// Run добавляет задачу в очередь и возвращает каналы для результата и ошибки
func Run(expression string) int {
	previousId++
	id := previousId

	expr := models.Expression{
		ID:     id,
		Expr:   expression,
		Status: "pending",
		Result: float64(0),
	}
	expressions[id] = expr
	updateExpressionStatus(id, "in_progress", float64(0))
	evaluateExpression(expr)
	return id
}

func updateExpressionStatus(id int, status string, result float64) {
	expressionsMu.Lock()
	defer expressionsMu.Unlock()
	if expr, ok := expressions[id]; ok {
		expr.Status = status
		expr.Result = result
		expressions[id] = expr
	}
}

func GetExpressions() []models.Expression {
	expressionsMu.Lock()
	defer expressionsMu.Unlock()

	var exprList []models.Expression
	for _, expr := range expressions {
		exprList = append(exprList, expr)
	}

	return exprList
}
func GetExpression(id int) models.Expression {
	expressionsMu.Lock()
	defer expressionsMu.Unlock()

	for _, expr := range expressions {
		if expr.ID == id {
			return expr
		}
	}

	return models.Expression{}
}
