package internal

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"tests-econtrols-supervisor/internal/entities"
)

var mqttAddressTemplate = "tcp://%s:%d"

type Worker struct {
	username string
	client   mqtt.Client
}

func initWorker(username string, config entities.SupervisorConfig) Worker {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf(mqttAddressTemplate, config.BrokerHost, config.BrokerPort)).SetUsername(username)

	client := mqtt.NewClient(opts)

	return Worker{username: username, client: client}
}

func work() {

}
