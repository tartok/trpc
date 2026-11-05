package trpc

import "fmt"

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Log     any    `json:"log"`
}

func (r RpcError) Error() string {
	return fmt.Sprintf("%s", r.Message)
}
