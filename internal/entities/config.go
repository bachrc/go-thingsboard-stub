package entities

import (
	"encoding/json"
	"io/ioutil"
)

type SupervisorConfig struct {
	Topics struct {
		Publish struct {
			Telemetry   string `json:"telemetry"`
			Attributes  string `json:"attributes"`
			RPCResponse string `json:"rpc_response"`
		} `json:"publish"`
		Subscribe struct {
			RPCRequests string `json:"rpc_requests"`
		} `json:"subscribe"`
		Regex struct {
			RPCRequests string `json:"rpc_requests"`
		} `json:"regex"`
	} `json:"topics"`
}

func GetConfig() (SupervisorConfig, error) {
	file, _ := ioutil.ReadFile("resources/supervisor_options.json")

	config := SupervisorConfig{}

	err := json.Unmarshal([]byte(file), &config)

	return config, err
}
