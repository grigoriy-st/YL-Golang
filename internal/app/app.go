package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/grigoriy-st/YL-Golang/pkg/calculator"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Структура выражения
type Expression struct {
	id     int     `json:"id"`
	exp    string  `json:"exp"`
	status string  `json:status`
	result float64 `json:result`
}

// Буфер задач
type SeqTasksBuffer struct {
	m         sync.Mutex
	buffer    []Expression
	idCounter int
}

// Возврат и удаление задачи
func (s *SeqTasksBuffer) popTask() (Expression, error) {
	s.m.Lock()
	defer s.m.Unlock()

	bufLenght := len(s.buffer)
	if bufLenght > 0 {
		last_exp := s.buffer[bufLenght-1]
		s.buffer = s.buffer[:bufLenght-1]
		return last_exp, nil
	}
	return Expression{}, fmt.Errorf("Error in pop task")
}

// Добавление новой задачи в буфер
func (s *SeqTasksBuffer) appendTask(task string) {
	s.m.Lock()
	defer s.m.Unlock()

	s.buffer = append(s.buffer, Expression{s.getIdForTask(), task, "Proccesed", 0.0})
}

func (s *SeqTasksBuffer) getIdForTask() int {
	s.m.Lock()
	defer s.m.Unlock()

	s.idCounter++
	return s.idCounter
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	buffer := SeqTasksBuffer{}
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		go func() {
			reader := bufio.NewReader(os.Stdin)

			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("failed to read expression from console")
			}
			// убираем пробелы, чтобы оставить только вычислемое выражение
			text = strings.TrimSpace(text)
			// выходим, если ввели команду "exit"
			if text == "exit" {
				log.Println("aplication was successfully closed")
				return
			}
			buffer.appendTask(text)
		}()
		//вычисляем выражение
		go func() {
			exp, err := buffer.popTask()
			if exp.status != "Proccesed" || err != nil {
				fmt.Errorf("Error in pop task")
			}
			result, err := calculator.Calc(exp.exp)
			if err != nil {
				log.Println(exp.exp, " calculator failed wit error: ", err)
			} else {
				log.Println(exp.exp, "=", result)
			}
		}()
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result"`
}

type Error struct {
	Error string `json:"error"`
}

// Обработчик выражений.
// Перенаправляет выражение в функцию, которая его вычсиляет
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultChan := make(chan *Response)
	errorChan := make(chan *Error)

	go func() {
		result, err := calculator.Calc(request.Expression)
		fmt.Println(result, err)
		if err != nil {
			if errors.Is(err, calculator.ErrDivisionByZero) {
				errorChan <- &Error{Error: "division by zero"}
			} else {
				errorChan <- &Error{Error: fmt.Sprintf("%v", err.Error())}
			}
			return
		}

		response := &Response{Result: fmt.Sprintf("%f", result)}
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

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
