package models

import (
	"fmt"
	"sync"
)

// Структура выражения
type Expression struct {
	id     int     `json:"id"`
	exp    string  `json:"exp"`
	status string  `json:status`
	result float64 `json:result`
}

func (e *Expression) NewExpression(exp string) Expression {
	return Expression{exp: exp}
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

// Получение уникального идентификатора
func (s *SeqTasksBuffer) getIdForTask() int {
	s.m.Lock()
	defer s.m.Unlock()

	s.idCounter++
	return s.idCounter
}
