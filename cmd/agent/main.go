package main

import (
	"log"

	"github.com/Dided08/Calculator/internal"
)

func main() {
	agent := application.NewAgent()
	log.Println("Starting Agent...")
	agent.Run()
}
