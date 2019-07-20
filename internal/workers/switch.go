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
	Client               *mqtt.Client
	Value                bool   `json:"value"`
	AttributeName        string `json:"attributeName"`
	GetValueMethod       string `json:"getValueMethod"`
	SetValueMethod       string `json:"setValueMethod"`
	getValueEventChannel *chan entities.RawRequest
	setValueEventChannel *chan entities.RawRequest
}

func (s *Switch) answerGetValue(request entities.RawRequest) {
	payload := make(map[string]bool)
	payload[s.AttributeName] = s.Value

	value := []byte(strconv.FormatBool(s.Value))
	(*s.Client).Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 1, false, value)
}

func (s *Switch) answerSetValue(request entities.RawRequest) {
	var setValueRequest entities.SetSwitchRequest
	unmarshalError := json.Unmarshal(request.Payload, &setValueRequest)

	log.Printf("The received request is %+v", setValueRequest)
	if unmarshalError != nil {
		log.Fatalf("Unparseable set temperature request : %b", request.Payload)
	}

	s.Value = setValueRequest.Params
	log.Printf("[SWIT : %10s] Set value to %s", s.AttributeName, setValueRequest.Params)

	(*s.Client).Publish(fmt.Sprintf(config.Topics.Publish.RPCResponse, request.RequestId), 1, false, request.Payload)
	s.sendValue()
}

func (s *Switch) sendValue() {
	payload := make(map[string]bool)

	payload[s.AttributeName] = s.Value

	messageToSend, _ := json.Marshal(payload)
	client := *s.Client
	token := client.Publish(config.Topics.Publish.Telemetry, 2, false, messageToSend)

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Printf("[SWIT : %10s] Message sent : %s", s.AttributeName, messageToSend)
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

func (s *Switch) SetupEventChannels(getValueEventChannel *chan entities.RawRequest, setValueEventChannel *chan entities.RawRequest) {
	s.getValueEventChannel = getValueEventChannel
	s.setValueEventChannel = setValueEventChannel
}
