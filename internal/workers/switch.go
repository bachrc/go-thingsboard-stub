package workers

import (
	"encoding/json"
	"fmt"
	"github.com/bachrc/thingsboard-stub/internal/entities"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

var config, _ = entities.GetConfig()

type Switch struct {
	Value                bool `json:"value"`
	Client               *mqtt.Client
	AttributeName        string `json:"attributeName"`
	GetValueMethod       string `json:"getValueMethod"`
	SetValueMethod       string `json:"setValueMethod"`
	getValueEventChannel *chan entities.RPCRequest
	setValueEventChannel *chan entities.RPCRequest
}

func (s *Switch) answerGetValue(request entities.RPCRequest) {
	payload := make(map[string]bool)
	payload[s.AttributeName] = s.Value

	response := entities.GetSwitchValue{
		Method: s.GetValueMethod,
		Params: entities.Params{
			Value: s.Value,
		},
	}

	messageToSend, _ := json.Marshal(response)
	(*s.Client).Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 2, false, messageToSend)
}

func (s *Switch) answerSetValue(request entities.RPCRequest) {
	newValue, err := strconv.ParseBool(request.Params)

	if err == nil {
		s.Value = newValue
		response := entities.SetSwitchValue{
			Method: "setValue",
			Params: newValue,
		}

		messageToSend, _ := json.Marshal(response)

		client := *s.Client
		client.Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 2, false, messageToSend)
	} else {
		log.Fatalf("Invalid switch value given : %s", request.Params)
	}

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
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case getValue := <-*s.getValueEventChannel:
			s.answerGetValue(getValue)
			break
		case setValue := <-*s.setValueEventChannel:
			s.answerSetValue(setValue)
			break
		case <-ticker.C:
			s.sendValue()
		}
	}
}

func (s *Switch) SetupEventChannels(getValueEventChannel *chan entities.RPCRequest, setValueEventChannel *chan entities.RPCRequest) {
	s.getValueEventChannel = getValueEventChannel
	s.setValueEventChannel = setValueEventChannel
}
