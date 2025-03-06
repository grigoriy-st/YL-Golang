package orchestrator

import (
	"calc/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
)

type Orchestrator struct {
	expressions map[int]models.Expression
	tasks       chan models.Task
	mutex       sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[int]models.Expression),
		tasks:       make(chan models.Task, 100),
	}
}

// Определение приоритета операции
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "(":
		return 0
	}
	return -1
}

func infixToRPN(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var output []string
	var operators []string
	i := 0
	for i < len(expression) {
		if expression[i] >= '0' && expression[i] <= '9' {
			num := ""
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				num += string(expression[i])
				i++
			}
			output = append(output, num)
			continue
		} else if expression[i] == '(' {
			operators = append(operators, string(expression[i]))
		} else if expression[i] == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(string(expression[i])) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, string(expression[i]))
		}
		i++
	}

	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// Добавление вырадение в оркестратор
func (o *Orchestrator) AddExpression(c echo.Context) error {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "Invalid JSON"})
	}

	exprID := rand.Intn(1000000) // генерация уникального идентификатора выражения
	o.mutex.Lock()
	o.expressions[exprID] = models.Expression{Id: exprID, Status: "processing"}
	o.mutex.Unlock()

	postfix, err := infixToRPN(req.Expression)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "Ошибка преобразования выражения"})
	}

	var arg1, arg2 float64
	for _, token := range postfix {
		if token == "+" || token == "-" || token == "*" || token == "/" {
			task := models.Task{
				Id:             exprID,
				Arg1:           arg1,
				Arg2:           arg2,
				Operation:      token,
				Operation_time: 10 * time.Millisecond,
			}
			o.tasks <- task
		} else {
			if num, err := strconv.ParseFloat(token, 64); err == nil {
				if arg1 == 0 {
					arg1 = num
				} else {
					arg2 = num
				}
			}
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": exprID})
}

func (o *Orchestrator) GetExpressionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"expressions": "Not found"})
	}
	o.mutex.Lock()
	expr, exists := o.expressions[id]
	o.mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{"expressions": "Not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"expression": expr})
}

func (o *Orchestrator) GetExpressions(c echo.Context) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if len(o.expressions) == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"expressions": "Not Found"})
	}

	var expressionsList []map[string]interface{}
	for _, exp := range o.expressions {
		expressionsList = append(expressionsList, map[string]interface{}{
			"id":     exp.Id,
			"status": exp.Status,
			"result": exp.Result,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"expressions": expressionsList})
}

func (o *Orchestrator) GetTask(c echo.Context) error {
	select {
	case task := <-o.tasks:
		return c.JSON(http.StatusOK, echo.Map{"task": task})
	default:
		return c.JSON(http.StatusNotFound, echo.Map{"error": "No tasks available"})
	}
}

// Получение результата
func (o *Orchestrator) ReceiveResult(c echo.Context) error {
	var result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := c.Bind(&result); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "Invalid JSON"})
	}

	o.mutex.Lock()
	expr, exists := o.expressions[result.ID]
	if exists {
		expr.Status = "completed"
		expr.Result = result.Result
		o.expressions[result.ID] = expr
	}
	o.mutex.Unlock()

	return c.NoContent(http.StatusOK)
}

// Запуск сервера
func StartServer() {
	e := echo.New()
	orchestrator := NewOrchestrator()

	e.POST("/api/v1/calculate", orchestrator.AddExpression)
	e.GET("/api/v1/expressions", orchestrator.GetExpressions)
	e.GET("/api/v1/expressions/:id", orchestrator.GetExpressionByID)
	e.GET("/internal/task", orchestrator.GetTask)
	e.POST("/internal/task/result", orchestrator.ReceiveResult)

	log.Println("Оркестратор запущен на порту 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
