package trpc

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRpc(t *testing.T) {
	r, err := NewRpcClient(Option[ClientBuilder]().
		Url("localhost:123").
		Transport(TransportTcp),
	)
	if err != nil {
		t.Fatal(err)
	}
	//var p any
	//p = struct {
	//	P1 string `json:"p_1"`
	//}{P1: "param 1"}
	for {
		for i := range 3 {
			go func() {
				resp, err := r.TcpRequest("test", i, nil)
				fmt.Println(resp, err)
			}()
		}
		time.Sleep(time.Second)
	}
}
