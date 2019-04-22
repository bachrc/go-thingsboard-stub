package entities

type GetValue struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}
type Params struct {
	Value bool `json:"value"`
}

type Instruction struct {
	Method string `json:"method"`
}

type SetValue struct {
	Method string `json:"method"`
	Params bool   `json:"params"`
}

type CheckStatus struct {
	Method string `json:"method"`
	Params bool   `json:"params"`
}
