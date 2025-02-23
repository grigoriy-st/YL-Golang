// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [orchestrator|agent]")
		return
	}

	switch os.Args[1] {
	case "orchestrator":
		cmd := exec.Command("go", "run", "orchestrator.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // Это нужно для отделения процесса
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Orchestrator started...")
	case "agent":
		cmd := exec.Command("go", "run", "agent.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Agent started...")
	default:
		fmt.Println("Unknown argument:", os.Args[1])
	}
}
