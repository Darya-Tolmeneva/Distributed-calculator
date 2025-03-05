package transport

import (
	"Distributed_calculator/internal/models"
	"Distributed_calculator/internal/services/orchestrator"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func PutExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var request models.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid request format"}`, http.StatusUnprocessableEntity)
		return
	}
	log.Println("Put expression to orchestrator")
	id := orchestrator.Run(request.Expression)
	response := map[string]int{
		"id": id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Failed to encode expressions"}`, http.StatusInternalServerError)
		return
	}

}
func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	exprList := orchestrator.GetExpressions()

	response := map[string][]models.Expression{
		"expressions": exprList,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Failed to encode expressions"}`, http.StatusInternalServerError)
		return
	}
}
func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/expressions/"):]
	id, _ := strconv.Atoi(idStr)
	expr := orchestrator.GetExpression(id)
	response := map[string]models.Expression{
		"expression": expr,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Failed to encode expressions"}`, http.StatusInternalServerError)
		return
	}
}

func HandleTaskResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result models.Result

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	orchestrator.GetTaskResult(result)
}
