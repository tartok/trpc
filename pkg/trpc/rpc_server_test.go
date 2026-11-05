package trpc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestRpcServer(t *testing.T) {
	srv, err := NewRpcServer(TransportTcp, ":123")
	if err != nil {
		t.Fatal(err)
	}
	srv.RegisterHandler(func(rpc *RpcServer, method string, params json.RawMessage) (any, any, bool, error) {
		fmt.Println(method, string(params))
		return nil, nil, true, nil
	})
	sigChan := time.After(10 * time.Second)
	//	sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	_ = srv.Close()
	time.Sleep(1 * time.Second)
}
