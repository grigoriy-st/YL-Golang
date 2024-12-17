package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculate(expression string) (string, error) {
	// Простой парсер для арифметических выражений
	// Здесь можно добавить более сложную логику для обработки выражений
	// В данном примере поддерживаются только простые операции +, -, *, /
	// и числа

	// Удаляем пробелы
	expression = strings.ReplaceAll(expression, " ", "")

	// Проверяем, что выражение состоит только из цифр и разрешённых операторов
	validExpression := regexp.MustCompile(`^[0-9+\-*/().]+$`)
	if !validExpression.MatchString(expression) {
		return "", fmt.Errorf("invalid expression")
	}

	// Используем встроенную функцию для вычисления выражения
	// В реальном приложении лучше использовать библиотеку для парсинга и вычисления
	result, err := eval(expression)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

// eval - простая функция для вычисления арифметического выражения
// В реальном приложении используйте более надежный парсер
func eval(expression string) (float64, error) {
	// Здесь можно использовать библиотеку для вычисления выражений
	// Для простоты, просто возвращаем 0
	return 0, nil // Замените на реальную логику вычисления
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result, err := calculate(req.Expression)
	if err != nil {
		if err.Error() == "invalid expression" {
			http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		} else {
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
