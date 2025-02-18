package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestCalcHandler_ParallelProcessing(t *testing.T) {
	t.Parallel() // Добавляем эту строку для параллельного выполнения теста

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(CalcHandler))
	defer server.Close()

	// Определяем тестовые выражения
	expressions := []string{
		// "1 + 1",
		// "2 * 3",
		// "4 / 2",
		// "5 - 3",
		"10 / 0", // Это выражение должно вызвать ошибку деления на ноль
	}

	var wg sync.WaitGroup
	results := make([]*Response, len(expressions))
	errors := make([]*Error, len(expressions))

	// Запус параллельных запросов
	for i, exp := range expressions {
		wg.Add(1)
		go func(i int, exp string) {
			defer wg.Done()

			// Создание JSON-запроса
			requestBody, _ := json.Marshal(Request{Expression: exp})
			resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Errorf("Failed to send request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var result Response
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Errorf("Failed to decode response: %v", err)
					return
				}
				results[i] = &result
			} else {
				var errResponse Error
				if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
					return
				}
				errors[i] = &errResponse
			}
		}(i, exp)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Проверяем результаты
	// expectedResults := []string{"2.000000", "6.000000", "2.000000", "2.000000", "division by zero"}
	expectedResults := []string{"division by zero"}
	for i, exp := range expressions {
		if errors[i] != nil {
			if exp == "10 / 0" && errors[i].Error != "division by zero" {
				t.Errorf("Expected error for expression '%s', got: %s", exp, errors[i].Error)
			}
		} else {
			if results[i].Result != expectedResults[i] {
				t.Errorf("Expected result for expression '%s', got: %s, want: %s", exp, results[i].Result, expectedResults[i])
			}
		}
	}
}
