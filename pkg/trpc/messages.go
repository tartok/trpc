package trpc

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Id = any
type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      Id              `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func (req *Request) Bytes() []byte {
	b, _ := json.Marshal(req)
	return b
}

func NewRequest(id Id, method string, params any) *Request {
	if id == nil {
		id = bson.NewObjectID()
	}
	p, _ := json.Marshal(params)
	r := Request{
		Jsonrpc: "2.0",
		Id:      id,
		Method:  method,
		Params:  p,
	}
	return &r
}

type Response struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      Id        `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *RpcError `json:"error,omitempty"`
}

func NewRpcSuccess(id Id, r any) Response {
	return Response{Jsonrpc: "2.0", Id: id, Result: r}
}
func NewRpcError(id Id, code int, message string, log any) Response {
	return Response{Jsonrpc: "2.0", Id: id, Error: &RpcError{Code: code, Message: message, Log: log}}
}
