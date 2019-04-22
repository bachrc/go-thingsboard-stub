package internal

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"tests-econtrols-supervisor/internal/entities"
)

var mqttAddressTemplate = "tcp://%s:%d"
var config, _ = entities.GetConfig()

type Worker struct {
	username string
	client   mqtt.Client
}

func (w *Worker) Work() {
	w.client.Connect()

	log.Println("Running...")
}

func (w *Worker) init(username string) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf(mqttAddressTemplate, config.BrokerHost, config.BrokerPort)).SetUsername(username)
	opts.OnConnect = w.onConnect

	w.client = mqtt.NewClient(opts)
}

func (w *Worker) onConnect(c mqtt.Client) {
	log.Print("I did connected well !")
	if token := c.Subscribe(config.Topics.DeviceAPI.Subscribe.RPCRequests, 2, w.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (w *Worker) onMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message from topic : %s", msg.Topic())
}

func InitWorker(username string) Worker {
	worker := Worker{}
	worker.init(username)

	return worker
}
