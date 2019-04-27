package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tests-econtrols-supervisor/internal"
	"tests-econtrols-supervisor/internal/entities"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	config, err := entities.GetConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("Bonjouuuur, voici le broker : %s", config.BrokerHost)

	worker := internal.InitWorker("v4b77JcaXdctUJtOTWoF", "valueSwitch1", "temperature1")

	go worker.Work()

	<-c
}
