package trpc

type RpcClient struct {
	opt ClientOptions
}

func NewRpcClient(o ...Lister[ClientOptions]) (*RpcClient, error) {
	var r RpcClient
	err := AddOptions(&r.opt, o...)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type ClientOptions struct {
	Id        any
	Url       string
	Transport Transport
	Auth      string
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
func (o *ClientBuilder) Transport(t Transport) *ClientBuilder {
	o.Opts = append(o.Opts, func(opts *ClientOptions) error {
		opts.Transport = t
		return nil
	})
	return o
}
