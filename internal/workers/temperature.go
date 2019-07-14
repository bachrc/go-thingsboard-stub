package workers

import (
	"../entities"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

type Temperature struct {
	Value          string
	Client         *mqtt.Client
	AttributeName  string
	GetValueMethod string
	SetValueMethod string
}

func (t *Temperature) answerGetValue(requestId string) {
	payload := make(map[string]string)
	payload[t.AttributeName] = t.Value

	response := entities.GetTemperatureValue{
		Method: t.GetValueMethod,
		Params: entities.ParamsTemperature{
			Value: t.Value,
		},
	}

	messageToSend, _ := json.Marshal(response)

	client := *t.Client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (t *Temperature) answerSetValue(message []byte, requestId string) {
	var receivedValue entities.SetTemperature
	_ = json.Unmarshal(message, &receivedValue)

	t.Value = receivedValue.Params
	response := entities.SetTemperature{
		Method: "setValue",
		Params: t.Value,
	}

	messageToSend, _ := json.Marshal(response)

	client := *t.Client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (t *Temperature) sendValue() {
	payload := make(map[string]string)

	payload[t.AttributeName] = t.Value
	log.Printf("Message sent : %+v", payload)

	messageToSend, _ := json.Marshal(payload)

	client := *t.Client
	client.Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)
}

func (t *Temperature) Work() {
	for range time.Tick(20 * time.Second) {
		t.sendValue()
	}
}

func (t *Temperature) HandleMessage(method string, requestId string, payload []byte) {
	switch method {
	case t.GetValueMethod:
		t.answerGetValue(requestId)
		break
	case t.SetValueMethod:
		t.answerSetValue(payload, requestId)
		break
	}
}
