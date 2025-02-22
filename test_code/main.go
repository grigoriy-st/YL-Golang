package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

var (
	expressions = make(map[string]*Expression)
	tasks       = make(chan Task, 100)
	mutex       = sync.Mutex{}
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	exprID := strconv.Itoa(rand.Int())
	expressions[exprID] = &Expression{ID: exprID, Status: "pending"}
	log.Printf("Received expression: %s", req.Expression)

	go parseExpression(exprID, req.Expression)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func parseExpression(id, expr string) {
	re := regexp.MustCompile(`([-+*/])`)
	parts := re.Split(expr, -1)
	operators := re.FindAllString(expr, -1)

	if len(parts) != 2 || len(operators) != 1 {
		mutex.Lock()
		expressions[id].Status = "error"
		mutex.Unlock()
		return
	}

	arg1, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	arg2, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		mutex.Lock()
		expressions[id].Status = "error"
		mutex.Unlock()
		return
	}

	operation := operators[0]
	task := Task{ID: id, Arg1: arg1, Arg2: arg2, Operation: operation, OperationTime: getOperationTime(operation)}
	tasks <- task
	mutex.Lock()
	expressions[id].Status = "processing"
	mutex.Unlock()
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-tasks:
		json.NewEncoder(w).Encode(map[string]Task{"task": task})
	default:
		http.Error(w, "No tasks", http.StatusNotFound)
	}
}

func postResultHandler(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := expressions[result.ID]
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	expr.Result = result.Result
	expr.Status = "completed"
	w.WriteHeader(http.StatusOK)
}

func getOperationTime(operation string) int {
	switch operation {
	case "+":
		return getEnvInt("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvInt("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	case "/":
		return getEnvInt("TIME_DIVISIONS_MS", 1000)
	default:
		return 1000
	}
}

func getEnvInt(name string, defaultVal int) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return defaultVal
	}
	return val
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")
	r.HandleFunc("/internal/task", getTaskHandler).Methods("GET")
	r.HandleFunc("/internal/task", postResultHandler).Methods("POST")

	http.Handle("/", r)
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
