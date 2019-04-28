package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tests-econtrols-supervisor/internal"
	"tests-econtrols-supervisor/internal/entities"
	"tests-econtrols-supervisor/internal/workers"
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

	switches := []workers.Switch{{
		AttributeName:  "valueSwitch1",
		Value:          true,
		GetValueMethod: "getSwitch1Value",
		SetValueMethod: "setSwitch1Value",
	}}

	temperatures := []workers.Temperature{{
		AttributeName:  "temperature1",
		Value:          "20.0",
		GetValueMethod: "getTemperature1",
		SetValueMethod: "setTemperature1",
	}}

	worker := internal.InitWorker("v4b77JcaXdctUJtOTWoF", switches, temperatures)

	go worker.Work()

	<-c
}
