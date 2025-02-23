// main.go
package main

import (
	"first/agent"
	"first/orchestrator"
	"os"
	"strconv"
)

func main() {
	// Получаем количество агентов из переменной окружения
	numAgents, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		numAgents = 10 // по умолчанию запускаем 10 агентов
	}

	// Запуск оркестратора
	go orchestrator.StartOrchestrator()

	// Запуск агентов
	for i := 0; i < numAgents; i++ {
		go agent.StartAgent(i)
	}

	// Блокируем основной поток
	select {}
}
