package models

// import (
// 	"fmt"
// 	"sync"
// )

// // Буфер выражений
// type SeqExpressionsBuffer struct {
// 	m         sync.Mutex
// 	buffer    []Expression
// 	idCounter int
// }

// // Возврат и удаление задачи
// func (s *SeqExpressionsBuffer) PopExpression() (Expression, error) {
// 	s.m.Lock()
// 	defer s.m.Unlock()

// 	bufLenght := len(s.buffer)
// 	if bufLenght > 0 {
// 		last_exp := s.buffer[bufLenght-1]
// 		s.buffer = s.buffer[:bufLenght-1]
// 		return last_exp, nil
// 	}
// 	return Expression{}, fmt.Errorf("Error in pop expression")
// }

// // Добавление новой задачи в буфер
// func (s *SeqExpressionsBuffer) AppendExpression(task string) {
// 	s.m.Lock()
// 	defer s.m.Unlock()
// 	fmt.Println("Добавление нового выражения в буфер")
// 	s.buffer = append(s.buffer, Expression{s.GetIdForExpression(), task, "Proccesed", 0.0})
// }

// // Получение уникального идентификатора выражения
// func (s *SeqExpressionsBuffer) GetIdForExpression() int {
// 	fmt.Println("Генерирование уникального идентификатора")
// 	s.idCounter++
// 	return s.idCounter
//  }
