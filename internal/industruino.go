package internal

import (
	"./entities"
	"./workers"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
)

var mqttAddressTemplate = "tcp://%s:%d"
var config, _ = entities.GetConfig()

type Device struct {
	username     string
	client       *mqtt.Client
	switches     []*workers.Switch
	temperatures []*workers.Temperature
}

func (w *Device) init(address string, port int, token string, switchesRef []*workers.Switch, temperaturesRef []*workers.Temperature) {
	w.username = token
	w.switches = switchesRef
	w.temperatures = temperaturesRef
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf(mqttAddressTemplate, address, port)).SetUsername(token)
	opts.OnConnect = w.onConnect

	client := mqtt.NewClient(opts)
	w.client = &client

	for _, theSwitch := range switchesRef {
		theSwitch.Client = &client
	}
	for _, theTemperature := range temperaturesRef {
		theTemperature.Client = &client
	}
}

func (w *Device) Work() {
	client := *w.client
	defer client.Disconnect(1)

	client.Connect()
	log.Println("Running...")

	for _, theSwitch := range w.switches {
		go (*theSwitch).Work()
	}

	for _, theTemperature := range w.temperatures {
		go (*theTemperature).Work()
	}
}

func (w *Device) onConnect(c mqtt.Client) {
	log.Println("I did connected well !")
	if token := c.Subscribe(config.Topics.Subscribe.RPCRequests, 2, w.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (w *Device) onMessage(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	log.Printf("Received message from topic : %s", msg.Topic())
	log.Printf("The message is : \n%s", payload)

	var received entities.Instruction
	_ = json.Unmarshal(payload, &received)

	requestId := w.getRequestId(msg.Topic())

	if received.Method == "checkStatus" {
		w.checkStatusHandler(requestId, payload)
	}

	for _, theSwitch := range w.switches {
		(*theSwitch).HandleMessage(received.Method, requestId, payload)
	}

	for _, theTemperature := range w.temperatures {
		(*theTemperature).HandleMessage(received.Method, requestId, payload)
	}

}

func (w *Device) getRequestId(topic string) string {
	r := regexp.MustCompile(config.Topics.Regex.RPCRequests)
	matches := r.FindStringSubmatch(topic)
	if len(matches) != 2 {
		panic("Invalid topic")
	}

	return matches[1]
}

func (w *Device) checkStatusHandler(requestId string, payload []byte) {
	client := *w.client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, payload)
}

func InitWorker(address string, port int, token string, switches []*workers.Switch, temperatures []*workers.Temperature) *Device {
	worker := new(Device)
	worker.init(address, port, token, switches, temperatures)

	return worker
}
