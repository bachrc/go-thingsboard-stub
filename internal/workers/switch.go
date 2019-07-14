package workers

import (
	"encoding/json"
	"fmt"
	"github.com/bachrc/thingsboard-stub/internal/entities"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

var config, _ = entities.GetConfig()

type Switch struct {
	Value          bool `json:"value"`
	Client         *mqtt.Client
	AttributeName  string `json:"attributeName"`
	GetValueMethod string `json:"getValueMethod"`
	SetValueMethod string `json:"setValueMethod"`
}

func (s *Switch) answerGetValue(requestId string) {
	payload := make(map[string]bool)
	payload[s.AttributeName] = s.Value

	response := entities.GetSwitchValue{
		Method: s.GetValueMethod,
		Params: entities.Params{
			Value: s.Value,
		},
	}

	messageToSend, _ := json.Marshal(response)
	client := *s.Client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (s *Switch) answerSetValue(message []byte, requestId string) {
	var receivedValue entities.SetSwitchValue
	_ = json.Unmarshal(message, &receivedValue)

	s.Value = receivedValue.Params
	response := entities.SetSwitchValue{
		Method: "setValue",
		Params: s.Value,
	}

	messageToSend, _ := json.Marshal(response)

	client := *s.Client
	client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, requestId), 2, false, messageToSend)
}

func (s *Switch) sendValue() {
	payload := make(map[string]bool)

	payload[s.AttributeName] = s.Value
	log.Printf("Message sent : %+v", payload)

	messageToSend, _ := json.Marshal(payload)
	client := *s.Client
	client.Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)
}

func (s *Switch) Work() {
	for range time.Tick(2 * time.Second) {
		s.sendValue()
	}
}

func (s *Switch) HandleMessage(method string, requestId string, payload []byte) {
	switch method {
	case s.GetValueMethod:
		s.answerGetValue(requestId)
		break
	case s.SetValueMethod:
		s.answerSetValue(payload, requestId)
		break
	}
}
