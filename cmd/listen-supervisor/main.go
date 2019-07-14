package main

import (
	"../../internal"
	"../../internal/entities"
	"../../internal/workers"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func altmain() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "gap",
			Value: 5,
			Usage: "Seconds gap between sending two values",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "start-temperature",
			Usage: "starts a temperature stub",
			Action: func(c *cli.Context) error {
				fmt.Println("added  ", c.Args().First())
				return nil
			},
		},
		{
			Name:  "switch",
			Usage: "starts a switch stub",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
