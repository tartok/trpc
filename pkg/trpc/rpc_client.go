package trpc

import (
	"encoding/json"
	"sync"
)

type queueChan struct {
	Req  *Request
	Resp chan []byte
}
type TransportFunc func(r *Request, open bool) (TransportClient, error)
type RpcClient struct {
	sync.Mutex
	opt       ClientOptions
	transport TransportFunc
}

func NewRpcClient(transport TransportFunc, o ...Lister[ClientOptions]) (*RpcClient, error) {
	var r = RpcClient{transport: transport}
	err := AddOptions(&r.opt, o...)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
func (r *RpcClient) Request(method string, params any, res any, customOptions ...Lister[ClientOptions]) (*Request, *Response, error) {
	options := r.opt
	err := AddOptions(&options, customOptions...)

	request := NewRequest(nil, method, params)
	tr, err := r.transport(request, true)
	if err != nil {
		return request, nil, err
	}
	data, err := tr.Request(request.Bytes())
	if err != nil {
		_, _ = r.transport(request, false)
		return request, nil, err
	}
	//	data, ok := <-tr.Receive()
	//	if !ok {
	//		_, _ = r.transport(request, false)
	//		return request, nil, errors.New("transport closed")
	//	}
	var result Response
	if res != nil {
		result.Result = res
	}
	err = json.Unmarshal(data, &result)
	return request, &result, err
}

type ClientOptions struct {
	Id  any
	Url string
	//Transport TransportServer
	Auth string
}
type ClientBuilder struct {
	OptionsBuilder[ClientOptions]
}

func (o *ClientBuilder) Id(id any) *ClientBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Id = id
		return nil
	})
	return o
}
func (o *ClientBuilder) Url(url string) *ClientBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Url = url
		return nil
	})
	return o
}
func (o *ClientBuilder) Auth(auth string) *ClientBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Auth = auth
		return nil
	})
	return o
}
