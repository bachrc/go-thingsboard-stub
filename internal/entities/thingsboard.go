package entities

type GetValue struct {
	EntityID string `json:"entityId"`
	OneWay   bool   `json:"oneWay"`
	Method   string `json:"method"`
	Params   struct {
		Value bool `json:"value"`
	} `json:"params"`
}
