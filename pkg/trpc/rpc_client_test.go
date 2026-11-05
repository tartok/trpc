package trpc

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewRpc(t *testing.T) {
	var tcp *TcpClient
	var m sync.Mutex
	defer func() {
		if tcp != nil {
			_ = tcp.Close()
		}
	}()
	r, err := NewRpcClient(func(r *Request) (t TransportClient, err error) {
		m.Lock()
		defer m.Unlock()
		if tcp == nil {
			tcp, err = NewTcpClient("localhost:123")
			if err != nil {
				return nil, err
			}
		}
		return tcp, err
	})
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
				var s string
				resp, err := r.Request("test", i, &s)
				fmt.Println(resp, err, s)
			}()
		}
		time.Sleep(time.Second)
	}
}
