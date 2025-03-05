package agent

import (
	"Distributed_calculator/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

func compute(task models.Task) (float64, error) {
	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
		time.Sleep(TIME_ADDITION_MS)
	case "-":
		result = task.Arg1 - task.Arg2
		time.Sleep(TIME_SUBTRACTION_MS)
	case "*":
		result = task.Arg1 * task.Arg2
		time.Sleep(TIME_MULTIPLICATIONS_MS)
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
		time.Sleep(TIME_DIVISIONS_MS)
	default:
		return 0, fmt.Errorf("invalid operation: %s", task.Operation)
	}
	return result, nil
}

func handleCompute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var operationTime time.Duration

	switch task.Operation {
	case "+":
		operationTime = TIME_ADDITION_MS
	case "-":
		operationTime = TIME_SUBTRACTION_MS
	case "*":
		operationTime = TIME_MULTIPLICATIONS_MS
	case "/":
		operationTime = TIME_DIVISIONS_MS
	}

	response := models.Task{
		ID:            task.ID,
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: operationTime,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	result, err := compute(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := sendResult(task.ID, result); err != nil {
		log.Printf("Failed to send result for task %d: %v", task.ID, err)
	} else {
		log.Printf("Result for task %d sent successfully: %f", task.ID, result)
	}
}

func sendResult(taskID int, result float64) error {
	resultData := models.Result{
		ID:     taskID,
		Result: result,
	}

	payload, err := json.Marshal(resultData)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit result: status %d", resp.StatusCode)
	}

	return nil
}

func Run(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/compute", handleCompute)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Agent started on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
}
