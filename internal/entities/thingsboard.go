package entities

type GetValue struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type GetTemperature struct {
	Method string            `json:"method"`
	Params ParamsTemperature `json:"params"`
}

type Params struct {
	Value bool `json:"value"`
}

type ParamsTemperature struct {
	Value float64 `json:"value"`
}

type Instruction struct {
	Method string `json:"method"`
}

type SetValue struct {
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
