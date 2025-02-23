// agent.go
package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID            int           `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

func StartAgent(workerID int) {
	for {
		task, err := getTaskFromOrchestrator()
		if err != nil {
			log.Printf("Agent %d: Error getting task: %s", workerID, err)
			time.Sleep(1 * time.Second) // Задержка, если задач нет
			continue
		}

		log.Printf("Agent %d: Processing task %d", workerID, task.ID)

		// Выполнение вычислений
		result := processTask(task)

		// Отправка результата в оркестратор
		err = submitTaskResultToOrchestrator(task.ID, result)
		if err != nil {
			log.Printf("Agent %d: Error submitting result: %s", workerID, err)
		}
	}
}

func getTaskFromOrchestrator() (Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return Task{}, fmt.Errorf("failed to get task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Task{}, fmt.Errorf("no tasks available, status: %d", resp.StatusCode)
	}

	var taskResp struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return Task{}, fmt.Errorf("failed to decode task response: %v", err)
	}

	return taskResp.Task, nil
}

func processTask(task Task) float64 {
	var result float64
	switch task.Operation {
	case "add":
		result = task.Arg1 + task.Arg2
	case "subtract":
		result = task.Arg1 - task.Arg2
	case "multiply":
		result = task.Arg1 * task.Arg2
	case "divide":
		if task.Arg2 != 0 {
			result = task.Arg1 / task.Arg2
		} else {
			log.Printf("Error: division by zero in task %d", task.ID)
			result = 0 // Неопределённый результат
		}
	}
	time.Sleep(task.OperationTime) // Эмуляция времени выполнения
	return result
}

func submitTaskResultToOrchestrator(taskID int, result float64) error {
	req := struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}{
		ID:     taskID,
		Result: result,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/internal/task/submit", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to submit result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit result, status: %d", resp.StatusCode)
	}

	return nil
}
