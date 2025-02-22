package agent

type Agent interface {
}

// Запускает COMPUTING_POWER горутин
func (a *Agent) StartRoutins() {

}

// Запрос задач у оркестратора
func (a *Agent) GetTasks() {

}

// Поиск свободного вычислителя
func (a *Agent) CheckFreeCalc() {

}

// Передача результата выражения оркестратору
func (a *Agent) SendResult() {

}
