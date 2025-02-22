package agent

import (
	"bytes"
	"calc/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type GoroutineState struct {
	Id   int  `json:"id"`
	Busy bool `json:"busy"`
}

type Agent struct {
	mu              sync.Mutex
	goroutineStates []GoroutineState
	tasksBuffer     models.SeqTasksBuffer
}

// Запускает COMPUTING_POWER горутин
func (a *Agent) StartRoutins(COMPUTING_POWER int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.goroutineStates = make([]GoroutineState, COMPUTING_POWER)

	for i := 0; i < COMPUTING_POWER; i++ {
		a.goroutineStates[i] = GoroutineState{Id: i, Busy: false}
		go a.worker(i)
	}
}

func (a *Agent) worker(id int) {
	a.mu.Lock()
	a.goroutineStates[id].Busy = true
	a.mu.Unlock()

	// task = a.GetTasks() // Взятие задачи на выполнение
	// result, err := функция_вычисления_выражения(task) и возврат в переменную
	// обработка ошибки

	a.mu.Lock()
	a.goroutineStates[id].Busy = false
	a.mu.Unlock()

}

// Запрос задачи у оркестратора
func (a *Agent) GetTask() models.Task {
	url := "http://internal/tasks"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return models.Task{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status")
		return models.Task{}
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return models.Task{}
	}

	return task
}

// Поиск свободного вычислителя
func (a *Agent) FindFreeCalc() []int {
	a.mu.Lock()
	defer a.mu.Unlock()

	freeGoroutines := []int{}
	for _, state := range a.goroutineStates {
		if !state.Busy {
			freeGoroutines = append(freeGoroutines, state.Id)
		}
	}
	return freeGoroutines
}

// Передача результата выражения оркестратору
func (a *Agent) SendResult(result float64) {
	url := "http://internal/task"
	jsonData := []byte(fmt.Sprintf(`{"result": "%s"}`, result))

	// Отправка POST-запроса
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}
