package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tartok/trpc/pkg/trpc"
)

func main() {
	typ := "tcp"
	//typ = "http"
	var r *trpc.RpcClient
	var err error
	switch typ {
	case "tcp":
		r, err = tcp()
		if err != nil {
			panic(err)
		}
	case "http":
		r, err = http()
		if err != nil {
			panic(err)
		}
	}
	var c int64
	//t := time.Now()
	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		cc := atomic.LoadInt64(&c)
	//		d := time.Since(t)
	//		fmt.Println(float64(cc) * float64(time.Second) / float64(d))
	//	}
	//}()
	wg := sync.WaitGroup{}
	for range 1 {
		wg.Go(func() {
			for {
				//<-w
				atomic.AddInt64(&c, 1)
				var s string
				req, resp, err := r.Request("test", nil, &s)
				if err != nil {
					fmt.Println(err)

				} else {
					_, _ = req, resp
					fmt.Println(req, resp, err, s)
				}
				fmt.Println(atomic.LoadInt64(&c))
				time.Sleep(time.Second)
			}
		})
	}
	wg.Wait()
}
