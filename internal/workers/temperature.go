package workers

import (
	"encoding/json"
	"fmt"
	"github.com/bachrc/thingsboard-stub/internal/entities"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

type Temperature struct {
	Client               *mqtt.Client
	Value                string
	AttributeName        string
	GetValueMethod       string
	SetValueMethod       string
	getValueEventChannel *chan entities.RPCRequest
	setValueEventChannel *chan entities.RPCRequest
}

func (t *Temperature) answerGetValue(request entities.RPCRequest) {
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
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 2, false, messageToSend)
}

func (t *Temperature) answerSetValue(request entities.RPCRequest) {
	t.Value = request.Params
	response := entities.SetTemperature{
		Method: t.SetValueMethod,
		Params: t.Value,
	}

	messageToSend, _ := json.Marshal(response)

	client := *t.Client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 2, false, messageToSend)
}

func (t *Temperature) sendValue() {
	payload := make(map[string]string)

	payload[t.AttributeName] = t.Value
	log.Printf("Message sent : %+v", payload)

	messageToSend, _ := json.Marshal(payload)
	log.Printf("Sending this json : %s to this topic : %s", messageToSend, config.Topics.Publish.Telemetry)
	client := *t.Client
	token := client.Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (t *Temperature) Work() {
	ticker := time.NewTicker(7 * time.Second)
	for {
		select {
		case getValue := <-*t.getValueEventChannel:
			t.answerGetValue(getValue)
			break
		case setValue := <-*t.setValueEventChannel:
			t.answerSetValue(setValue)
			break
		case <-ticker.C:
			t.sendValue()
		}
	}
}

func (t *Temperature) SetupEventChannels(getValueEventChannel *chan entities.RPCRequest, setValueEventChannel *chan entities.RPCRequest) {
	t.getValueEventChannel = getValueEventChannel
	t.setValueEventChannel = setValueEventChannel
}
