package main

import (
	"../../internal"
	"../../internal/utils"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var gap int
	var address string
	var port int
	var token string
	var switches string
	var temperatures string

	app := cli.NewApp()

	app.Name = "ecstub"
	app.Usage = "Stub for the eControls Supervisor run by Thingsboard"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "gap, g",
			Value:       5,
			Usage:       "Seconds gap between sending two values",
			Destination: &gap,
		},
		cli.StringFlag{
			Name:        "address, a",
			Usage:       "The address of the broker",
			Destination: &address,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       1883,
			Usage:       "Port of the MQTT broker",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "token",
			Usage:       "The token (identifier) of the device on thingsboard",
			Destination: &token,
		},
		cli.StringFlag{
			Name:        "switches, s",
			Usage:       "Path to the switches definition file",
			Value:       "resources/default_switches.json",
			Destination: &switches,
		},
		cli.StringFlag{
			Name:        "temperatures, t",
			Usage:       "Path to the temperatures definition file",
			Value:       "resources/default_temperatures.json",
			Destination: &temperatures,
		},
	}

	app.Action = func(c *cli.Context) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		industruino := startApplication(gap, address, port, token, switches, temperatures)
		industruino.Work()

		<-ch

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startApplication(gap int, address string, port int, token string, switchesDefPath string, temperaturesDefPath string) *internal.Device {
	switches, _ := utils.ParseSwitchesDefinition(switchesDefPath)
	temperatures, _ := utils.ParseTemperaturesDefinition(temperaturesDefPath)

	worker := internal.InitWorker(gap, address, port, token, switches, temperatures)

	return worker
}
