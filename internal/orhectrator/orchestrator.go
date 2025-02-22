package orchestrator

type Orchestrator interface {
}

// Проверка буфера на свободные задачи
func (o *Orchestrator) CheckBuffer() {

}

// Дробление выражение на задачи
func (o *Orchestrator) ParseExpIntoTasks(exp string) {

}

// Выдача результатов
func (o *Orchestrator) GiveResult() {

}
