package entities

type RPCRequest struct {
	Method    string `json:"method"`
	RequestId string
}

type RawRequest struct {
	RPCRequest
	Payload []byte
}

type SetTemperatureRequest struct {
	RPCRequest
	Value string `json:"value"`
}

type SetSwitchRequest struct {
	RPCRequest
	Value bool `json:"value"`
}
