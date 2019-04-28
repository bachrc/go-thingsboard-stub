package entities

type GetSwitchValue struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type GetTemperatureValue struct {
	Method string            `json:"method"`
	Params ParamsTemperature `json:"params"`
}

type Params struct {
	Value bool `json:"value"`
}

type ParamsTemperature struct {
	Value string `json:"value"`
}

type Instruction struct {
	Method string `json:"method"`
}

type SetSwitchValue struct {
	Method string `json:"method"`
	Params bool   `json:"params"`
}

type SetTemperature struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type CheckStatus struct {
	Method string `json:"method"`
	Params bool   `json:"params"`
}
