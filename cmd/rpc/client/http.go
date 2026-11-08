package main

import "github.com/tartok/trpc/pkg/trpc"

func http() (*trpc.RpcClient, error) {
	r, err := trpc.NewRpcClient(func(r *trpc.Request, open bool) (trpc.TransportClient, error) {
		if open {
			rc, err := trpc.NewHttpClient("http://localhost:8080/rpc")
			return rc, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
