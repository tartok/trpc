package trpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestRpcServer(t *testing.T) {
	tcp, err := NewTcpServer(":123")
	if err != nil {
		t.Fatal(err)
	}
	defer tcp.Close()
	go tcp.Listen()
	srv, err := NewRpcServer(tcp)
	if err != nil {
		t.Fatal(err)
	}
	var c int32
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("c:", atomic.LoadInt32(&c))
		}
	}()

	srv.RegisterHandler(func(rpc *RpcServer, method string, params json.RawMessage) (any, any, bool, error) {
		//fmt.Println(method, string(params))
		//time.Sleep(time.Microsecond)
		if rand.Intn(10) > 7 {
			time.Sleep(time.Microsecond)
		}
		//for i := 0; i < rand.Intn(1000); i++ {
		//
		//}
		atomic.AddInt32(&c, 1)
		return "OK", nil, true, nil
	})
	sigChan := time.After(1000 * time.Second)
	//	sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	_ = srv.Close()
	time.Sleep(1 * time.Second)
}
