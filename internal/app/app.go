package app

import (
	"bufio"
	"calc/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	buffer := models.SeqTasksBuffer{}
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

			buffer.AppendTask(text)
		}()
		//вычисляем выражение
		go func() {
			exp, err := buffer.PopTask()
			if exp.Status != "Proccesed" || err != nil {
				fmt.Errorf("Error in pop task")
			}
			result, err := calculator.Calc(&exp)
			if err != nil {
				log.Println(exp.Exp, " calculator failed wit error: ", err)
			} else {
				log.Println(exp.Exp, "=", result)
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

	buffer := models.SeqTasksBuffer{}

	go func() {
		buffer.AppendTask(request.Expression)
		exp, _ := buffer.PopTask()

		result, err := calculator.Calc(&exp)

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
