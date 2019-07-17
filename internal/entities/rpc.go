package entities

type RPCRequest struct {
	Method    string `json:"method"`
	Params    string `json:"params"`
	RequestId string
}
