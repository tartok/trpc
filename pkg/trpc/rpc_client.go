package trpc

import (
	"encoding/json"
	"sync"
)

type queueChan struct {
	Req  *Request
	Resp chan []byte
}

type RpcClient struct {
	sync.Mutex
	queue     []*queueChan
	opt       ClientOptions
	transport func(r *Request) (TransportClient, error)
	wait      chan struct{}
}

func NewRpcClient(transport func(r *Request) (TransportClient, error), o ...Lister[ClientOptions]) (*RpcClient, error) {
	var r = RpcClient{transport: transport, wait: make(chan struct{}, 100)}
	err := AddOptions(&r.opt, o...)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			<-r.wait
			tr, err := r.transport(nil)
			if err != nil {
				continue
			}
			buff := <-tr.Receive()
			var result RespId
			err = json.Unmarshal(buff, &result)
			if err != nil {
				continue
			}
			r.Lock()
			for _, q := range r.queue {
				if q.Req != nil && q.Req.Id == result.Id {
					q.Req = nil
					q.Resp <- buff
				}
			}
			r.Unlock()
		}
	}()
	return &r, nil
}
func (r *RpcClient) send(req *Request) <-chan []byte {
	r.Lock()
	defer r.Unlock()
	for _, q := range r.queue {
		if q.Req == nil {
			q.Req = req
			return q.Resp
		}
	}
	q := queueChan{
		Req:  req,
		Resp: make(chan []byte),
	}
	r.queue = append(r.queue, &q)
	return q.Resp
}
func (r *RpcClient) Request(method string, params any, res any, customOptions ...Lister[ClientOptions]) (*Response, error) {
	options := r.opt
	err := AddOptions(&options, customOptions...)

	request := NewRequest(nil, method, params)
	c := r.send(request)

	tr, err := r.transport(request)
	if err != nil {
		return nil, err
	}
	r.wait <- struct{}{}
	_, err = tr.Write(request.Bytes())
	if err != nil {
		return nil, err
	}
	data := <-c
	var result Response
	if res != nil {
		result.Result = res
	}
	err = json.Unmarshal(data, &result)
	return &result, nil
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
