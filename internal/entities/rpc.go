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
	Params string `json:"params"`
}

type SetSwitchRequest struct {
	RPCRequest
	Params bool `json:"params"`
}
