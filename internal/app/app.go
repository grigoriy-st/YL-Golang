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
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
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
			return nil
		}
		//вычисляем выражение
		result, err := calculator.Calc(text)
		if err != nil {
			log.Println(text, " calculator failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
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

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calculator.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculator.ErrInvalidExpression) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrInvalidExpression) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrDivisionByZero) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrIncorrectSeqOfParenthese) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrDiffNumberOfBrackets) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrConvertingNumberToFloatType) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrTwoOperatorsInRow) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrTwoOperandsInRow) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrExpStartsWithOperator) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else if errors.Is(err, calculator.ErrExpEndsWithOperator) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else {
			fmt.Fprintf(w, "Unknown err")
		}

	} else {
		fmt.Fprintf(w, "result: %f", result)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
