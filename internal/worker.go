package internal

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
	"tests-econtrols-supervisor/internal/entities"
	"tests-econtrols-supervisor/internal/workers"
)

var mqttAddressTemplate = "tcp://%s:%d"
var config, _ = entities.GetConfig()

type Industruino struct {
	username     string
	client       mqtt.Client
	switches     []workers.Switch
	temperatures []workers.Temperature
}

func (w *Industruino) Work() {
	defer w.client.Disconnect(1)

	w.client.Connect()
	log.Println("Running...")
	for _, theSwitch := range w.switches {
		go theSwitch.Work()
	}

	for _, theTemperature := range w.temperatures {
		go theTemperature.Work()
	}
}

func (w *Industruino) init(username string, switches []workers.Switch, temperatures []workers.Temperature) {
	w.username = username
	w.switches = switches
	w.temperatures = temperatures
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf(mqttAddressTemplate, config.BrokerHost, config.BrokerPort)).SetUsername(username)
	opts.OnConnect = w.onConnect

	w.client = mqtt.NewClient(opts)

	for _, theSwitch := range w.switches {
		theSwitch.Client = w.client
	}
	for _, theTemperature := range w.temperatures {
		theTemperature.Client = w.client
	}
}

func (w *Industruino) onConnect(c mqtt.Client) {
	log.Print("I did connected well !")
	if token := c.Subscribe(config.Topics.Subscribe.RPCRequests, 2, w.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (w *Industruino) onMessage(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	log.Printf("Received message from topic : %s", msg.Topic())
	log.Printf("The message is : \n%s", payload)

	var received entities.Instruction
	_ = json.Unmarshal(payload, &received)

	requestId := w.getRequestId(msg.Topic())

	if received.Method == "checkStatus" {
		w.checkStatusHandler(client, requestId, payload)
	}

	for _, theSwitch := range w.switches {
		theSwitch.HandleMessage(received.Method, requestId, payload)
	}

	for _, theTemperature := range w.temperatures {
		theTemperature.HandleMessage(received.Method, requestId, payload)
	}

}

func (w *Industruino) getRequestId(topic string) string {
	r := regexp.MustCompile(config.Topics.Regex.RPCRequests)
	matches := r.FindStringSubmatch(topic)
	if len(matches) != 2 {
		panic("Invalid topic")
	}

	return matches[1]
}

func (w *Industruino) checkStatusHandler(client mqtt.Client, requestId string, payload []byte) {
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, payload)
}

func InitWorker(username string, switches []workers.Switch, temperatures []workers.Temperature) *Industruino {
	worker := new(Industruino)
	worker.init(username, switches, temperatures)

	return worker
}
