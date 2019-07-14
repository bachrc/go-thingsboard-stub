package utils

import (
	"encoding/json"
	"github.com/bachrc/thingsboard-stub/internal/workers"
	"io/ioutil"
)

func ParseSwitchesDefinition(path string) ([]*workers.Switch, error) {
	file, _ := ioutil.ReadFile(path)

	var switches []*workers.Switch

	err := json.Unmarshal([]byte(file), &switches)

	return switches, err
}

func ParseTemperaturesDefinition(path string) ([]*workers.Temperature, error) {
	file, _ := ioutil.ReadFile(path)

	var temperatures []*workers.Temperature

	err := json.Unmarshal([]byte(file), &temperatures)

	return temperatures, err
}
