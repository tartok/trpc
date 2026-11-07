package trpc

import (
	"fmt"
	"sync"
	"sync/atomic"
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
	r, err := NewRpcClient(func(r *Request, open bool) (t TransportClient, err error) {
		m.Lock()
		defer m.Unlock()
		if open {
			if tcp == nil {
				tcp, err = NewTcpClient("localhost:123")
				if err != nil {
					return nil, err
				}
			}
			return tcp, nil
		}
		if tcp != nil {
			_ = tcp.Close()
			tcp = nil
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	//var p any
	//p = struct {
	//	P1 string `json:"p_1"`
	//}{P1: "param 1"}
	var c, e, p int32
	go func() {
		for {
			time.Sleep(1 * time.Second)
			cc := atomic.LoadInt32(&c)
			ee := atomic.LoadInt32(&e)
			fmt.Println("c:", cc, float64(ee)*100/float64(cc), cc-p)
			p = cc
		}
	}()
	w := make(chan struct{})
	for i := range 20 {
		go func() {
			for {
				//<-w
				var s string
				req, resp, err := r.Request("test", i, &s)
				if err != nil {
					fmt.Println(err)
					continue
				}
				atomic.AddInt32(&c, 1)
				if req.Id != resp.Id {
					atomic.AddInt32(&e, 1)
					fmt.Println(req, resp, err, s)
				}
			}
		}()
	}
	for {
		close(w)
		w = make(chan struct{})
		time.Sleep(time.Second)
	}
}
