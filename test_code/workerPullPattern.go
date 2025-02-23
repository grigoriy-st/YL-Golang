package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Id             int           `json:"id"`
	Arg1           float64       `json:"arg1"`
	Arg2           float64       `json:"arg2"`
	Operation      string        `json:"operation"`
	Operation_time time.Duration `json:"operation_time"`
}

type Expression struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

const (
	TIME_ADDITION_MS        = 100 * time.Millisecond
	TIME_SUBTRACTION_MS     = 100 * time.Millisecond
	TIME_MULTIPLICATIONS_MS = 100 * time.Millisecond
	TIME_DIVISIONS_MS       = 100 * time.Millisecond
)

// Запрос задачи у оркестратора
func GetTask() {
	url := "http://internal/tasks"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status")
		fmt.Errorf("%v", err)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		fmt.Println("Error decoding JSON:", err)
		fmt.Errorf("%v", err)
	}

	tasks <- task
}

func sendResult(result Expression) {
	url := "http://internal/task"

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Result sent successfully!")
	} else {
		fmt.Printf("Failed to send result, status code: %d\n", resp.StatusCode)
	}
}

func worker(id int, tasks <-chan Task, results chan<- float64) {
	for task := range tasks {
		fmt.Printf("worker %d started task %d\n", id, task.Id)
		time.Sleep(task.Operation_time)

		var result float64
		switch task.Operation {
		case "addition":
			result = task.Arg1 + task.Arg2
		case "subtraction":
			result = task.Arg1 - task.Arg2
		case "multiplication":
			result = task.Arg1 * task.Arg2
		case "division":
			if task.Arg2 != 0 {
				result = task.Arg1 / task.Arg2
			} else {
				fmt.Printf("worker %d: division by zero in task %d\n", id, task.Id)
				result = 0 // Или обработка ошибки
			}
		default:
			fmt.Printf("worker %d: unknown operation %s in task %d\n", id, task.Operation, task.Id)
			continue
		}

		fmt.Printf("worker %d finished task %d with result: %f\n", id, task.Id, result)
		results <- result

		// Отправка результата по HTTP POST запросу
		expression := Expression{
			Id:     task.Id,
			Result: result,
		}
		sendResult(expression)
	}
}

func main() {
	tasks := make(chan Task, 100)
	results := make(chan float64, 100)

	// Запускаем 3 worker'а
	for w := 1; w <= 3; w++ {
		go worker(w, tasks, results)
	}

	// Горутина для получения результатов
	go func() {
		for result := range results {
			fmt.Println("Result:", result)
		}
	}()

	// Бесконечный цикл для добавления новых задач
	go func() {
		for i := 1; ; i++ {
			// Пример создания задачи
			task := Task{
				Id:             i,
				Arg1:           float64(i),
				Arg2:           float64(i + 1),
				Operation:      "addition", // Здесь можно менять операции
				Operation_time: TIME_ADDITION_MS,
			}
			tasks <- task
			time.Sleep(500 * time.Millisecond) // Задержка между задачами
		}
	}()

	// Даем программе работать в течение 10 секунд, затем закрываем каналы
	time.Sleep(10 * time.Second)
	close(tasks)
	close(results)
}
