package main

import (
	"sync"

	"github.com/tartok/trpc/pkg/trpc"
)

func tcp() (*trpc.RpcClient, error) {
	var tcp *trpc.TcpClient
	var m sync.Mutex
	defer func() {
		if tcp != nil {
			_ = tcp.Close()
		}
	}()
	r, err := trpc.NewRpcClient(func(r *trpc.Request, open bool) (t trpc.TransportClient, err error) {
		m.Lock()
		defer m.Unlock()
		if open {
			if tcp == nil {
				tcp, err = trpc.NewTcpClient("localhost:123")
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
		return nil, err
	}
	return r, nil
}
