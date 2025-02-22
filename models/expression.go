package models

import (
	"fmt"
	"sync"
)

// Структура выражения
type Expression struct {
	Id     int     `json:"id"`
	Exp    string  `json:"exp"`
	Status string  `json:status`
	Result float64 `json:result`
}

// Создание нового выражения
func (e *Expression) NewExpression(exp string) Expression {
	return Expression{Exp: exp}
}

// Буфер задач
type SeqTasksBuffer struct {
	m         sync.Mutex
	buffer    []Expression
	idCounter int
}

// Возврат и удаление задачи
func (s *SeqTasksBuffer) PopTask() (Expression, error) {
	s.m.Lock()
	defer s.m.Unlock()

	bufLenght := len(s.buffer)
	if bufLenght > 0 {
		last_exp := s.buffer[bufLenght-1]
		s.buffer = s.buffer[:bufLenght-1]
		last_exp.Status = "Finished"
		return last_exp, nil
	}
	return Expression{}, fmt.Errorf("Error in pop task")
}

// Добавление новой задачи в буфер
func (s *SeqTasksBuffer) AppendTask(task string) {
	s.m.Lock()
	defer s.m.Unlock()
	fmt.Println("Добавление новой задачи в буфер")
	s.buffer = append(s.buffer, Expression{s.GetIdForTask(), task, "Being processed", 0.0})
}

// Получение уникального идентификатора для задачи
func (s *SeqTasksBuffer) GetIdForTask() int {
	fmt.Println("Генерирование уникального идентиффикатора")
	s.idCounter++
	return s.idCounter
}
