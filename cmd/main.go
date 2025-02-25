package main

import (
	"calc/internal/agent"
	orchestrator "calc/internal/orchectrator"

	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	serverURL := "http://localhost:8080"
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || computingPower <= 0 {
		computingPower = 10 // Значение по умолчанию
	}
	wg.Add(2)

	// Запуск оркестратора
	go func() {
		defer wg.Done()
		log.Println("Запуск оркестратора...")
		// orchestrator.StartServer()
		orchestrator.StartServer()
	}()

	// Запуск агента
	go func() {
		defer wg.Done()
		log.Println("Запуск агента...")
		agent.StartServer(serverURL, computingPower)
	}()

	wg.Wait()
}
