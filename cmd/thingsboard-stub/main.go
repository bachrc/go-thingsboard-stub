package main

import (
	"github.com/bachrc/thingsboard-stub/internal"
	"github.com/bachrc/thingsboard-stub/internal/utils"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var address string
	var port int
	var token string
	var switches string
	var temperatures string

	app := cli.NewApp()

	app.Name = "thingsboard-stub"
	app.Usage = "Stub for Thingsboard"

	app.Flags = []cli.Flag{
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

		if address == "" {
			panic("Please provide the address of the Thingsboard instance")
		}

		if token == "" {
			panic("Please provide the Thingsboard Acces Token of the device to stub")
		}

		industruino := createDevice(address, port, token, switches, temperatures)
		defer industruino.Stop()
		industruino.Work()

		<-ch
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createDevice(address string, port int, token string, switchesDefPath string, temperaturesDefPath string) *internal.Device {
	switches, _ := utils.ParseSwitchesDefinition(switchesDefPath)
	temperatures, _ := utils.ParseTemperaturesDefinition(temperaturesDefPath)

	worker := internal.InitWorker(address, port, token, switches, temperatures)

	return worker
}
