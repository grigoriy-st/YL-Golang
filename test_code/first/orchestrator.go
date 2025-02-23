// orchestrator.go
package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Структуры для выражений и задач
type Expression struct {
	ID         int     `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result,omitempty"`
}

type Task struct {
	ID            int           `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

var (
	expressions = make(map[int]Expression)
	taskQueue   = make(chan Task, 100)
	mu          sync.Mutex
)

func StartOrchestrator() {
	// Настройки времени выполнения операций из окружения
	operationTimes := map[string]time.Duration{
		"addition":       parseDuration(os.Getenv("TIME_ADDITION_MS")),
		"subtraction":    parseDuration(os.Getenv("TIME_SUBTRACTION_MS")),
		"multiplication": parseDuration(os.Getenv("TIME_MULTIPLICATIONS_MS")),
		"division":       parseDuration(os.Getenv("TIME_DIVISIONS_MS")),
	}

	http.HandleFunc("/api/v1/calculate", calculateExpression)
	http.HandleFunc("/api/v1/expressions", getExpressions)
	http.HandleFunc("/api/v1/expressions/", getExpressionByID)
	http.HandleFunc("/internal/task", getTask)
	http.HandleFunc("/internal/task/submit", submitTaskResult)

	// Запуск HTTP сервера оркестратора
	fmt.Println("Orchestrator is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseDuration(ms string) time.Duration {
	parsed, err := strconv.Atoi(ms)
	if err != nil {
		return 0
	}
	return time.Millisecond * time.Duration(parsed)
}

// HTTP Handlers для Оркестратора

func calculateExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	// Генерация уникального ID
	mu.Lock()
	id := rand.Intn(10000)
	mu.Unlock()

	// Разбиение выражения на задачи
	tasks, err := parseExpressionToTasks(req.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	// Сохранение выражения в мапе
	expressions[id] = Expression{
		ID:         id,
		Expression: req.Expression,
		Status:     "received",
	}

	// Отправка задач в очередь
	for _, task := range tasks {
		taskQueue <- task
	}

	// Ответ с ID выражения
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func parseExpressionToTasks(expression string) ([]Task, error) {
	tokens := strings.Fields(expression)
	tasks := []Task{}

	for i := 0; i < len(tokens)-2; i += 2 {
		arg1, _ := strconv.ParseFloat(tokens[i], 64)
		arg2, _ := strconv.ParseFloat(tokens[i+2], 64)
		operation := tokens[i+1]
		tasks = append(tasks, Task{
			ID:            rand.Intn(10000),
			Arg1:          arg1,
			Arg2:          arg2,
			Operation:     operation,
			OperationTime: time.Millisecond * 100, // Время выполнения операции
		})
	}

	return tasks, nil
}

func getExpressions(w http.ResponseWriter, r *http.Request) {
	var resp []Expression
	for _, exp := range expressions {
		resp = append(resp, exp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Expression{"expressions": resp})
}

func getExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	exprID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expr, ok := expressions[exprID]
	if !ok {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]Expression{"expression": expr})
}

func getTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-taskQueue:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]Task{"task": task})
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func submitTaskResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	// Обновляем статус выражения с результатом
	for i := range expressions {
		if expressions[i].ID == req.ID {
			expressions[i].Result = req.Result
			expressions[i].Status = "completed"
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
