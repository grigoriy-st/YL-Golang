package orchestrator

import (
	"calc/models"
	"regexp"
	"strconv"
)

type Orchestrator struct {
}

// Проверка буфера на свободные задачи
func (o *Orchestrator) CheckBuffer() {

}

// Дробление выражение на задачи
func (o *Orchestrator) ParseExpIntoTasks(exp string) {
	tokens := tokenize(exp)
	rpn := shuntingYard(tokens)
	return convertToTasks(rpn), nil
}

// Токенизация выражения
func tokenize(expr string) []string {
	re := regexp.MustCompile(`\s*([+*/()-]|[0-9]*\.?[0-9]+)\s*`)
	return re.FindAllString(expr, -1)
}

// Алгоритм Шунтинг-Ярд для преобразования в ОПН
func shuntingYard(tokens []string) []string {
	var output []string
	var stack []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Удаляем "("
		} else {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

// Проверка, является ли токен числом
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// Преобразование ОПН в задачи
func convertToTasks(rpn []string) []models.Task {
	var stack []models.Task
	var tasks []models.Task
	taskBuffer := NewSeqTasksBuffer(100)

	for _, token := range rpn {
		if isNumber(token) {
			value, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, models.Task{Arg1: value})
		} else {
			arg2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			arg1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			task := models.Task{
				Id:        taskBuffer.GetIdForTask(),
				Arg1:      arg1.Arg1,
				Arg2:      arg2.Arg1,
				Operation: token,
			}
			tasks = append(tasks, task)
			taskBuffer.AppendTask(task)
		}
	}

	return tasks
}

// Выдача результатов
func (o *Orchestrator) GiveResult() {

}
