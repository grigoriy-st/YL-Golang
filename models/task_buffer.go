package models

// import (
// 	"fmt"
// 	"sync"
// )

// // Буфер задач
// type SeqTasksBuffer struct {
// 	m         sync.Mutex
// 	ch        chan Task
// 	idCounter int
// 	closed    bool
// }

// // Создание экзеляра буфера задач
// func NewSeqTasksBuffer(size int) *SeqTasksBuffer {
// 	return &SeqTasksBuffer{
// 		ch:     make(chan Task, size),
// 		closed: false,
// 	}
// }

// // Возврат и удаление задачи
// func (s *SeqTasksBuffer) PopTask() (Task, error) {
// 	s.m.Lock()
// 	defer s.m.Unlock()

// 	var taskForEx Task
// 	taskForEx, ok := <-s.ch
// 	if ok {
// 		return taskForEx, nil
// 	}
// 	return Task{}, fmt.Errorf("Error in pop task")
// }

// // Добавление новой задачи в буфер
// func (s *SeqTasksBuffer) AppendTask(task Task) {
// 	s.m.Lock()
// 	defer s.m.Unlock()
// 	fmt.Println("Добавление новой задачи в буфер")
// 	s.ch <- task
// }

// // Получение уникального идентификатора
// func (s *SeqTasksBuffer) GetIdForTask() int {
// 	fmt.Println("Генерирование уникального идентиффикатора")
// 	s.idCounter++
// 	return s.idCounter
// }
