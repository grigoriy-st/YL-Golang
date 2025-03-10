package agent

import (
	"bytes"
	"calc/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Agent struct {
	serverURL     string
	mutex         sync.Mutex
	currentResult float64 // Добавлено для хранения текущего результата
}

func NewAgent(serverURL string, computingPower int) *Agent {
	agent := &Agent{
		serverURL:     serverURL,
		currentResult: 0, // Инициализация текущего результата
	}
	for i := 0; i < computingPower; i++ {
		go agent.worker()
	}
	return agent
}

func (a *Agent) worker() {
	for {
		resp, err := http.Get(a.serverURL + "/internal/task")
		if err != nil {
			log.Println("Ошибка запроса задачи:", err)
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			log.Println("Нет задач, ждем...")
			time.Sleep(time.Second)
			continue
		}

		var task struct {
			Task models.Task `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			log.Println("Ошибка декодирования задачи:", err)
			continue
		}

		resp.Body.Close()

		log.Printf("Выполняю задачу ID %d: %.2f %s %.2f\n", task.Task.Id, task.Task.Arg1, task.Task.Operation, task.Task.Arg2)
		time.Sleep(task.Task.Operation_time)

		// Используем текущий результат для вычисления
		var result float64
		if task.Task.Id == 1 { // Если это первая задача, используем Arg1
			result = task.Task.Arg1
		} else {
			result = a.currentResult // Используем текущий результат для последующих задач
		}

		switch task.Task.Operation {
		case "+":
			result += task.Task.Arg2
		case "-":
			result -= task.Task.Arg2
		case "*":
			result *= task.Task.Arg2
		case "/":
			if task.Task.Arg2 != 0 {
				result /= task.Task.Arg2
			} else {
				log.Printf("Ошибка: деление на ноль в задаче ID %d", task.Task.Id)
				continue
			}
		}

		a.currentResult = result // Сохраняем текущий результат

		resultData, _ := json.Marshal(map[string]interface{}{"id": task.Task.Id, "result": result})
		_, err = http.Post(a.serverURL+"/internal/task/result", "application/json", bytes.NewBuffer(resultData))
		if err != nil {
			log.Printf("Ошибка отправки результата задачи ID %d: %v", task.Task.Id, err)
		}
	}
}

func StartServer(serverURL string, computingPower int) {
	NewAgent(serverURL, computingPower)
	select {} // Запуск в фоновом режиме
}
