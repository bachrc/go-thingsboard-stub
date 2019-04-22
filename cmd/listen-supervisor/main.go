package main

import (
	"fmt"
	"os"
	"tests-econtrols-supervisor/internal/entities"
)

func main() {
	config, err := entities.GetConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Bonjouuuur, voici le broker : %s", config.BrokerHost)
}
