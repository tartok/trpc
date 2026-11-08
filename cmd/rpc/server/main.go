package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tartok/trpc/pkg/trpc"
)

func main() {
	var transport trpc.TransportServer
	switch "http" {
	case "http":
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		httpSrv, err := trpc.NewHttpServer()
		if err != nil {
			panic(err)
		}
		h := http.NewServeMux()
		h.HandleFunc("/rpc", httpSrv.Handle)
		go func() {
			err := http.Serve(l, h)
			if err != nil {
				fmt.Println(err)
			}
		}()
		transport = httpSrv
	case "tcp":
		tcp, err := trpc.NewTcpServer(":123")
		if err != nil {
			panic(err)
		}
		defer tcp.Close()
		go tcp.Listen()
		transport = tcp
	}
	srv, err := trpc.NewRpcServer(transport)
	if err != nil {
		panic(err)
	}
	defer srv.Close()
	srv.RegisterHandler(func(rpc *trpc.RpcServer, method string, params json.RawMessage) (any, any, bool, error) {
		return "OK", nil, true, nil
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

}
