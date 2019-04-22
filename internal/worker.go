package internal

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
	"tests-econtrols-supervisor/internal/entities"
	"time"
)

var mqttAddressTemplate = "tcp://%s:%d"
var config, _ = entities.GetConfig()

type Worker struct {
	username      string
	attributeKey  string
	client        mqtt.Client
	booleanSwitch bool
}

func (w *Worker) Work() {
	defer w.client.Disconnect(1)

	w.client.Connect()
	log.Println("Running...")
	for range time.Tick(5 * time.Second) {
		w.sendValue()
	}
}

func (w *Worker) init(username string, attributeKey string) {
	w.username = username
	w.attributeKey = attributeKey
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf(mqttAddressTemplate, config.BrokerHost, config.BrokerPort)).SetUsername(username)
	opts.OnConnect = w.onConnect

	w.client = mqtt.NewClient(opts)
}

func (w *Worker) onConnect(c mqtt.Client) {
	log.Print("I did connected well !")
	if token := c.Subscribe(config.Topics.Subscribe.RPCRequests, 2, w.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (w *Worker) onMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message from topic : %s", msg.Topic())
	log.Printf("The message is : \n%s", msg.Payload())
	var received entities.Instruction
	_ = json.Unmarshal(msg.Payload(), &received)

	requestId := w.getRequestId(msg.Topic())

	switch received.Method {
	case "getValue":
		w.getValueHandler(client, requestId, msg.Payload())
		break
	case "setValue":
		w.setValueHandler(client, requestId, msg.Payload())
		break
	case "checkStatus":
		w.checkStatusHandler(client, requestId, msg.Payload())
		break
	}
}

func (w *Worker) AnswerToGetValue(client mqtt.Client, topic string, operation string) {
	requestId := w.getRequestId(topic)

	response := entities.SetValue{
		Method: operation,
		Params: true,
	}

	message, _ := json.Marshal(response)

	log.Printf("And you answer this object : %s", message)

	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, message)
}

func (w *Worker) getRequestId(topic string) string {
	r := regexp.MustCompile(config.Topics.Regex.RPCRequests)
	matches := r.FindStringSubmatch(topic)
	if len(matches) != 2 {
		panic("Invalid topic")
	}

	return matches[1]
}

func (w *Worker) getValueHandler(client mqtt.Client, requestId string, message []byte) {
	var received entities.GetValue
	_ = json.Unmarshal(message, &received)

	response := entities.GetValue{
		Method: "getValue",
		Params: entities.Params{
			Value: w.booleanSwitch,
		},
	}

	messageToSend, _ := json.Marshal(response)

	log.Printf("And you answer this object : %s", messageToSend)

	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (w *Worker) setValueHandler(client mqtt.Client, requestId string, message []byte) {
	var receivedValue entities.SetValue
	_ = json.Unmarshal(message, &receivedValue)

	w.booleanSwitch = receivedValue.Params
	response := entities.SetValue{
		Method: "setValue",
		Params: w.booleanSwitch,
	}

	messageToSend, _ := json.Marshal(response)

	log.Printf("And you answer this object : %s", messageToSend)

	topicToAnswer := fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId)
	client.Publish(topicToAnswer, 2, false, messageToSend)
}

func (w *Worker) checkStatusHandler(client mqtt.Client, requestId string, payload []byte) {
	response := entities.CheckStatus{
		Method: "checkStatus",
		Params: w.booleanSwitch,
	}

	messageToSend, _ := json.Marshal(response)

	log.Printf("And you answer this object : %s", messageToSend)

	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (w *Worker) sendValue() {
	payload := make(map[string]bool)

	payload[w.attributeKey] = w.booleanSwitch
	log.Printf("Message sent : %+v", payload)

	messageToSend, _ := json.Marshal(payload)

	w.client.Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)
}

func InitWorker(username string, attributeKey string) *Worker {
	worker := new(Worker)
	worker.init(username, attributeKey)

	return worker
}
