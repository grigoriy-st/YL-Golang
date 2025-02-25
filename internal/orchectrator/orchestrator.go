package orchestrator

import (
	"calc/models"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Проверка буфера на свободные задачи
func (o *Orchestrator) CheckBuffer() {

}

type Orchestrator struct {
	expressions map[int]models.Expression
	tasks       chan models.Task
	mutex       sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[int]models.Expression),
		tasks:       make(chan models.Task, 100),
	}
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "(":
		return 0
	}
	return -1
}

func infixToRPN(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var output []string
	var operators []string
	i := 0
	for i < len(expression) {
		if expression[i] >= '0' && expression[i] <= '9' {
			num := ""
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				num += string(expression[i])
				i++
			}
			output = append(output, num)
			continue
		} else if expression[i] == '(' {
			operators = append(operators, string(expression[i]))
		} else if expression[i] == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(string(expression[i])) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, string(expression[i]))
		}
		i++
	}

	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func (o *Orchestrator) AddExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	exprID := rand.Intn(1000000)
	o.mutex.Lock()
	o.expressions[exprID] = models.Expression{Id: exprID, Status: "processing"}
	o.mutex.Unlock()

	postfix, err := infixToRPN(req.Expression)
	if err != nil {
		http.Error(w, "Ошибка преобразования выражения", http.StatusUnprocessableEntity)
		return
	}

	// Разбиваем выражение и создаем задачи
	var arg1, arg2 float64
	for _, token := range postfix {
		if token == "+" || token == "-" || token == "*" || token == "/" {
			// Создание задачи для операции
			task := models.Task{
				Id:             exprID,
				Arg1:           arg1,
				Arg2:           arg2,
				Operation:      token,
				Operation_time: 10 * time.Millisecond,
			}
			o.tasks <- task
		} else {
			// Преобразуем токен в число
			if num, err := strconv.ParseFloat(token, 64); err == nil {
				// Присваиваем значение аргумента
				if arg1 == 0 {
					arg1 = num
				} else {
					arg2 = num
				}
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": exprID})
}

func (o *Orchestrator) GetExpressions(w http.ResponseWriter, r *http.Request) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": o.expressions})
}

func (o *Orchestrator) GetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-o.tasks:
		json.NewEncoder(w).Encode(map[string]models.Task{"task": task})
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func (o *Orchestrator) ReceiveResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	o.mutex.Lock()
	expr, exists := o.expressions[result.ID]
	if exists {
		expr.Status = "completed"
		expr.Result = result.Result
		o.expressions[result.ID] = expr
	}
	o.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func StartServer() {
	orchestrator := NewOrchestrator()
	http.HandleFunc("/api/v1/calculate", orchestrator.AddExpression)
	http.HandleFunc("/api/v1/expressions", orchestrator.GetExpressions)
	http.HandleFunc("/internal/task", orchestrator.GetTask)
	http.HandleFunc("/internal/task/result", orchestrator.ReceiveResult)

	log.Println("Оркестратор запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
