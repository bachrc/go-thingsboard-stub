package internal

import (
	"encoding/json"
	"fmt"
	"github.com/bachrc/thingsboard-stub/internal/entities"
	"github.com/bachrc/thingsboard-stub/internal/workers"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
)

var mqttAddressTemplate = "tcp://%s:%d"
var config, _ = entities.GetConfig()

type Device struct {
	username     string
	client       mqtt.Client
	operations   map[string]*chan entities.RawRequest
	switches     []*workers.Switch
	temperatures []*workers.Temperature
}

func (w *Device) init(address string, port int, token string, switchesRef []*workers.Switch, temperaturesRef []*workers.Temperature) {
	w.username = token
	w.switches = switchesRef
	w.temperatures = temperaturesRef
	mqttBrokerAddress := fmt.Sprintf(mqttAddressTemplate, address, port)
	log.Printf("The device is connecting to %s", mqttBrokerAddress)
	opts := mqtt.NewClientOptions().
		AddBroker(mqttBrokerAddress).
		SetUsername(token).
		SetAutoReconnect(true).
		SetPingTimeout(500).
		SetWriteTimeout(500).
		SetClientID("tb-stub")
	opts.OnConnect = (*w).onConnect

	w.client = mqtt.NewClient(opts)

	w.operations = make(map[string]*chan entities.RawRequest)

	for _, theSwitch := range switchesRef {
		_, okGet := w.operations[theSwitch.GetValueMethod]
		_, okSet := w.operations[theSwitch.SetValueMethod]

		if okGet || okSet {
			panic("Duplicate method names")
		}

		getValueEventChannel := make(chan entities.RawRequest)
		w.operations[theSwitch.GetValueMethod] = &getValueEventChannel
		setValueEventChannel := make(chan entities.RawRequest)
		w.operations[theSwitch.SetValueMethod] = &setValueEventChannel

		theSwitch.SetupEventChannels(&getValueEventChannel, &setValueEventChannel)
		theSwitch.Client = &w.client
	}
	for _, theTemperature := range temperaturesRef {
		_, okGet := w.operations[theTemperature.GetValueMethod]
		_, okSet := w.operations[theTemperature.SetValueMethod]

		if okGet || okSet {
			panic("Duplicate method names")
		}
		getValueEventChannel := make(chan entities.RawRequest)
		setValueEventChannel := make(chan entities.RawRequest)
		w.operations[theTemperature.GetValueMethod] = &getValueEventChannel
		w.operations[theTemperature.SetValueMethod] = &setValueEventChannel

		theTemperature.SetupEventChannels(&getValueEventChannel, &setValueEventChannel)
		theTemperature.Client = &w.client
	}
}

func (w *Device) Work() {
	w.client.Connect()
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
	rpcRequestsTopic := config.Topics.Subscribe.RPCRequests
	if token := c.Subscribe(rpcRequestsTopic, 2, w.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("Successfully subscribed to topic %s", rpcRequestsTopic)
	}
}

func (w *Device) onMessage(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	log.Printf("Received message %s from topic : %s", payload, msg.Topic())

	var received entities.RawRequest
	err := json.Unmarshal(payload, &received)

	if err != nil {
		log.Fatal("Error while parsing ")
	}

	received.RequestId = getRequestId(msg.Topic())
	received.Payload = msg.Payload()

	w.notifyEvent(received)
}

func (w *Device) notifyEvent(received entities.RawRequest) {
	eventChannel, present := w.operations[received.Method]
	if !present {
		log.Printf("Worker with method %s not found", received.Method)
		return
	}

	*eventChannel <- received
}

func (w *Device) checkStatusHandler(requestId string, payload []byte) {
	w.client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, payload)
}

func (w *Device) Stop() {
	log.Println("Disconnecting the client...")
	w.client.Disconnect(1)
}

func InitWorker(address string, port int, token string, switches []*workers.Switch, temperatures []*workers.Temperature) *Device {
	worker := new(Device)
	worker.init(address, port, token, switches, temperatures)

	return worker
}

func getRequestId(topic string) string {
	r := regexp.MustCompile(config.Topics.Regex.RPCRequests)
	matches := r.FindStringSubmatch(topic)
	if len(matches) != 2 {
		panic("Invalid topic")
	}

	return matches[1]
}
