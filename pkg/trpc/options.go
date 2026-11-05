package trpc

import (
	"fmt"
)

type Transport string

const (
	Http      Transport = "http"
	Websocket Transport = "websocket"
)

type Options struct {
	Url       string
	Transport Transport
}

type OptionsBuilder struct {
	Opts []func(*Options) error
}

func Option() *OptionsBuilder {
	return &OptionsBuilder{}
}

func (o *OptionsBuilder) SetUrl(url string) *OptionsBuilder {
	o.Opts = append(o.Opts, func(opts *Options) error {
		opts.Url = url
		return nil
	})
	return o
}
func (o *OptionsBuilder) SetTransport(t Transport) *OptionsBuilder {
	o.Opts = append(o.Opts, func(opts *Options) error {
		opts.Transport = t
		return nil
	})
	return o
}
func NewOptions(opts ...*OptionsBuilder) (*Options, error) {
	var o Options
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		for _, setArgs := range opt.Opts {
			if setArgs == nil {
				continue
			}

			if err := setArgs(&o); err != nil {
				return nil, err
			}
		}
	}
	return &o, nil
}
func test() {
	o := Option().SetTransport(Websocket).SetUrl("")
	fmt.Println(o)
}
