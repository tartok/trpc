package trpc

import (
	"encoding/json"
	"errors"
	"net"
	"sync"
)

type RpcHandler func(rpc *RpcServer, method string, params json.RawMessage) (any, any, bool, error)

type RpcServer struct {
	sync.Mutex
	opt        ServerOptions
	transports []TransportServer
	handlers   []RpcHandler
	listeners  []net.Listener
}

func NewRpcServer(transports []TransportServer, o ...Lister[ServerOptions]) (*RpcServer, error) {
	var r = RpcServer{transports: transports}
	err := AddOptions(&r.opt, o...)
	if err != nil {
		return nil, err
	}
	for _, transport := range transports {
		go func() {
			for {
				select {
				case msg, ok := <-transport.Get():
					if !ok {
						return
					}
					var req Request
					err := json.Unmarshal(msg.Request, &req)
					if err != nil {
						continue
					}
					h := false
					for _, handler := range r.handlers {
						if handler != nil {
							res, log, ok, err := handler(&r, req.Method, req.Params)
							if !ok {
								continue
							}
							var resp Response
							if err != nil {
								if errors.As(err, new(*RpcError)) {
									resp = NewRpcError(req.Id, err.(*RpcError).Code, err.(*RpcError).Message, log)
								} else {
									resp = NewRpcError(req.Id, -32603, err.Error(), log)
								}
							} else {
								resp = NewRpcSuccess(req.Id, res)
							}
							_, _ = msg.Writer.Write(resp.Bytes())
							h = true
							break
						}
					}
					if !h {
						resp := NewRpcError(req.Id, -32601, "Method not found", nil)
						_, _ = msg.Writer.Write(resp.Bytes())
					}
				}
			}
		}()
	}
	return &r, nil
}
func (r *RpcServer) Close() error {
	r.Lock()
	defer r.Unlock()
	for _, l := range r.listeners {
		_ = l.Close()
	}
	return nil
}
func (r *RpcServer) RegisterHandler(handler RpcHandler) {
	r.Lock()
	defer r.Unlock()
	r.handlers = append(r.handlers, handler)
}

type ServerOptions struct {
	//Url       string
	//Transport Transport
}
type ServerBuilder struct {
	OptionsBuilder[ClientOptions]
}

//func (o *ServerBuilder) Url(url string) *ServerBuilder {
//	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
//		opts.Url = url
//		return nil
//	})
//	return o
//}
//func (o *ServerBuilder) Transport(t Transport) *ServerBuilder {
//	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
//		opts.Transport = t
//		return nil
//	})
//	return o
//}
