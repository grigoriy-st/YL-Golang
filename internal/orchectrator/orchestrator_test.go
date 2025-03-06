package orchestrator_test

import (
	orchestrator "calc/internal/orchectrator"
	"calc/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo"
)

func TestOrchestratorParallel(t *testing.T) {
	e := echo.New()
	orch := orchestrator.NewOrchestrator()
	e.POST("/api/v1/calculate", orch.AddExpression)
	e.GET("/api/v1/expressions", orch.GetExpressions)
	e.GET("/api/v1/expressions/:id", orch.GetExpressionByID)
	e.GET("/internal/task", orch.GetTask)
	e.POST("/internal/task/result", orch.ReceiveResult)

	var wg sync.WaitGroup
	cases := []struct {
		name       string
		expression string
		statusCode int
	}{
		{"valid addition", "3+5", http.StatusCreated},
		{"valid multiplication", "2*8", http.StatusCreated},
		{"invalid expression", "3+/5", http.StatusUnprocessableEntity},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wg.Add(1)
			go func(tc struct {
				name       string
				expression string
				statusCode int
			}) {
				defer wg.Done()
				reqBody := `{"expression": "` + tc.expression + `"}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()

				e.ServeHTTP(respRec, req)

				if respRec.Code != tc.statusCode {
					t.Errorf("Expected status %d, got %d", tc.statusCode, respRec.Code)
				}
			}(tc)
		})
	}

	wg.Wait()
}

func TestOrchestratorTaskProcessing(t *testing.T) {
	e := echo.New()
	orch := orchestrator.NewOrchestrator()
	e.POST("/api/v1/calculate", orch.AddExpression)
	e.GET("/internal/task", orch.GetTask)
	e.POST("/internal/task/result", orch.ReceiveResult)

	var wg sync.WaitGroup
	reqBody := `{"expression": "7-2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	respRec := httptest.NewRecorder()
	e.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Fatalf("Failed to add expression: expected 201, got %d", respRec.Code)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Даем время оркестратору обработать выражение
		// time.Sleep(100 * time.Millisecond)
		req := httptest.NewRequest(http.MethodGet, "/internal/task", nil)
		respRec := httptest.NewRecorder()
		e.ServeHTTP(respRec, req)

		if respRec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", respRec.Code)
		}

		var taskResp struct {
			Task models.Task `json:"task"`
		}
		json.Unmarshal(respRec.Body.Bytes(), &taskResp)

		resultBody := `{"id": ` + strconv.Itoa(taskResp.Task.Id) + `, "result": 5}`
		resReq := httptest.NewRequest(http.MethodPost, "/internal/task/result", strings.NewReader(resultBody))
		resReq.Header.Set("Content-Type", "application/json")
		resRec := httptest.NewRecorder()
		e.ServeHTTP(resRec, resReq)

		if resRec.Code != http.StatusOK {
			t.Errorf("Expected status 200 on result submission, got %d", resRec.Code)
		}
	}()

	wg.Wait()
}
