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
	getValueEventChannel *chan entities.RawRequest
	setValueEventChannel *chan entities.RawRequest
}

func (t *Temperature) answerGetValue(request entities.RawRequest) {
	payload := make(map[string]string)
	payload[t.AttributeName] = t.Value

	value := []byte(t.Value)
	(*t.Client).Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 1, false, value)
}

func (t *Temperature) answerSetValue(request entities.RawRequest) {
	var setValueRequest entities.SetTemperatureRequest
	unmarshalError := json.Unmarshal(request.Payload, &setValueRequest)

	if unmarshalError != nil {
		log.Fatalf("Unparseable set temperature request : %b", request.Payload)
		return
	}

	t.Value = setValueRequest.Params

	(*t.Client).Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 1, false, request.Payload)
	t.sendValue()
}

func (t *Temperature) sendValue() {
	payload := make(map[string]string)

	payload[t.AttributeName] = t.Value

	messageToSend, _ := json.Marshal(payload)
	token := (*t.Client).Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Printf("[TEMP : %10s] Message sent : %s", t.AttributeName, messageToSend)
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

func (t *Temperature) SetupEventChannels(getValueEventChannel *chan entities.RawRequest, setValueEventChannel *chan entities.RawRequest) {
	t.getValueEventChannel = getValueEventChannel
	t.setValueEventChannel = setValueEventChannel
}
