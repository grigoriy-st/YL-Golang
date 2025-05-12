package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type requestPayload struct {
	Expression string `json:"expression"`
}

type responsePayload struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func TestCalculateHandler_Success(t *testing.T) {
	// Подготовка JSON-запроса
	body, _ := json.Marshal(requestPayload{Expression: "10 + 5 * 2"})
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Создание response recorder и вызов хендлера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalculateHandler)
	handler.ServeHTTP(rr, req)

	// Проверка кода ответа
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// Распарсить JSON-ответ
	var resp responsePayload
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Проверка результата
	expected := 20.0 // 10 + 5*2
	if resp.Result != expected {
		t.Errorf("expected result %.2f, got %.2f", expected, resp.Result)
	}
}

func TestCalculateHandler_BadExpression(t *testing.T) {
	body, _ := json.Marshal(requestPayload{Expression: "3 + "})
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalculateHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusOK {
		t.Errorf("expected 200 or 400 for bad input, got %d", rr.Code)
	}

	var resp responsePayload
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected an error message in response")
	}
}

func TestCalculateHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalculateHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", rr.Code)
	}
}
