package app

import (
	"bufio"
	"calc/api/handler"
	"calc/internal/config"
	"calc/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grigoriy-st/YL-Golang/pkg/calculator"
)

type Application struct {
	Config *config.Config
}

func New() *Application {
	return &Application{
		Config: config.ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение

func (a *Application) Run() error {
	buffer := models.SeqTasksBuffer{}
	taskChannel := make(chan string)              // Канал для передачи задач
	resultChannel := make(chan models.Expression) // Канал для передачи результатов

	// Горутина для обработки вычислений
	go func() {
		for {
			task := <-taskChannel // Получаем задачу из канала
			exp := models.Expression{Id: buffer.GetIdForTask(), Exp: task, Status: "Processed", Result: 0.0}
			result, err := calculator.Calc(&exp)
			if err != nil {
				log.Println(exp.Exp, "calculator failed with error:", err)
			} else {
				resultChannel <- models.Expression{Exp: exp.Exp, Result: result} // Отправляем результат в канал
			}
		}
	}()

	// Горутина для обработки результатов
	go func() {
		for result := range resultChannel {
			log.Println(result.Exp, "=", result.Result) // Выводим результат
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
			continue
		}

		// убираем пробелы, чтобы оставить только вычисляемое выражение
		text = strings.TrimSpace(text)

		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("application was successfully closed")
			close(resultChannel) // Закрываем канал результатов
			return nil
		}

		// Добавляем задачу в буфер и отправляем в канал
		buffer.AppendTask(text)
		taskChannel <- text // Отправляем задачу в канал
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", handler.CalcHandler)
	http.HandleFunc("/internal/tasks", handler.GetTaskHandler)
	return http.ListenAndServe(":"+a.Config.Addr, nil)
}
