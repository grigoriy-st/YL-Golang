package handler

import (
	"calc/models"
	"calc/pkg/calculator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Обработчик выражений.
// Перенаправляет выражение в функцию, которая его вычсиляет
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed in Server mode", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Начал обработку в CalcHandler")
	w.Header().Set("Content-Type", "application/json")
	request := new(models.Request)
	fmt.Println("Создал запрос")
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	fmt.Println("Декоирую ответ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultChan := make(chan *models.Response)
	errorChan := make(chan *models.Error)

	buffer := models.SeqTasksBuffer{}
	fmt.Println("Буфер создан")
	go func() {
		fmt.Println("Первая функция пошла")
		buffer.AppendTask(request.Expression)
		fmt.Println("Выражение добавлено")
		exp, err := buffer.PopTask()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		fmt.Println(exp)
		result, err := calculator.Calc(&exp)
		fmt.Println("Result:", result)
		if err != nil {
			if errors.Is(err, calculator.ErrDivisionByZero) {
				errorChan <- &models.Error{Error: "division by zero"}
			} else {
				errorChan <- &models.Error{Error: fmt.Sprintf("%v", err.Error())}
			}
			return
		}

		response := &models.Response{Result: fmt.Sprintf("%f", result)}
		resultChan <- response
	}()

	select {
	case response := <-resultChan:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case responseErr := <-errorChan:
		if responseErr.Error == "division by zero" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		json.NewEncoder(w).Encode(responseErr)
	}
}
