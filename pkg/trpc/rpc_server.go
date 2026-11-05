package trpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"sync"
)

type RpcHandler func(rpc *RpcServer, method string, params json.RawMessage) (any, any, bool, error)

type rpcMessage struct {
	request *Request
	w       func(p []byte) (n int, err error)
}

type RpcServer struct {
	sync.Mutex
	opt       ServerOptions
	inQueue   chan *rpcMessage
	handlers  []RpcHandler
	listeners []net.Listener
}

func NewRpcServer(transport Transport, url string, o ...Lister[ServerOptions]) (*RpcServer, error) {
	var r = RpcServer{inQueue: make(chan *rpcMessage, 100)}
	r.opt.Url = url
	r.opt.Transport = transport
	err := AddOptions(&r.opt, o...)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case msg := <-r.inQueue:
				for _, handler := range r.handlers {
					if handler != nil {
						res, log, ok, err := handler(&r, msg.request.Method, msg.request.Params)
						if !ok {
							continue
						}
						var resp Response
						if err != nil {
							if errors.As(err, new(*RpcError)) {
								resp = NewRpcError(msg.request.Id, err.(*RpcError).Code, err.(*RpcError).Message, log)
							} else {
								resp = NewRpcError(msg.request.Id, -32603, err.Error(), log)
							}
						} else {
							resp = NewRpcSuccess(msg.request.Id, res)
						}
						d := resp.Bytes()
						d = append(d, 0)
						_, _ = msg.w(d)
					}
				}
			}
		}
	}()
	switch r.opt.Transport {
	case TransportTcp:
		l, err := net.Listen("tcp", r.opt.Url)
		if err != nil {
			return nil, err
		}
		r.listeners = append(r.listeners, l)
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				go func(conn net.Conn) {
					reader := bufio.NewReader(conn)
					var m sync.Mutex
					for {
						data, err := reader.ReadBytes(0)
						if err != nil {
							return
						}
						var req Request
						err = json.Unmarshal(data[:len(data)-1], &req)
						if err != nil {
							continue
						}
						r.inQueue <- &rpcMessage{request: &req, w: func(p []byte) (n int, err error) {
							m.Lock()
							defer m.Unlock()
							return conn.Write(p)
						}}
					}
				}(conn)
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
	Url       string
	Transport Transport
}
type ServerBuilder struct {
	OptionsBuilder[ClientOptions]
}

func (o *ServerBuilder) Url(url string) *ServerBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Url = url
		return nil
	})
	return o
}
func (o *ServerBuilder) Transport(t Transport) *ServerBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Transport = t
		return nil
	})
	return o
}
