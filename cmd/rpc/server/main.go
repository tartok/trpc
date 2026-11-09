package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tartok/trpc/pkg/trpc"
)

func main() {
	var transports []trpc.TransportServer
	{
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
		transports = append(transports, httpSrv)
	}
	{
		var myKey = []byte{1, 2, 3}
		tcp, err := trpc.NewTcpServer(":123")
		if err != nil {
			panic(err)
		}
		defer tcp.Close()
		go tcp.Listen(func(r io.Reader) error {
			key := make([]byte, len(myKey))
			_, err := r.Read(key)
			if err != nil {
				return err
			}
			if bytes.Compare(key, myKey) != 0 {
				return errors.New("key not match")
			}
			return nil
		})
		transports = append(transports, tcp)
	}
	srv, err := trpc.NewRpcServer(transports)
	if err != nil {
		panic(err)
	}
	defer srv.Close()
	srv.RegisterHandler(func(rpc *trpc.RpcServer, method string, params json.RawMessage) (any, any, bool, error) {
		switch method {
		case "test":
			return "OK", nil, true, nil
		default:
			return nil, nil, false, nil
		}
	})
	srv.RegisterHandler(func(rpc *trpc.RpcServer, method string, params json.RawMessage) (any, any, bool, error) {
		switch method {
		case "test2":
			return "OK", nil, true, nil
		default:
			return nil, nil, false, nil
		}
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

}
