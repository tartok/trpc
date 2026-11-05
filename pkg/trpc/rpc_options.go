package trpc

type Transport string

const (
	TransportTcp       Transport = "tcp"
	TransportHttp      Transport = "http"
	TransportWebsocket Transport = "websocket"
)

type OptionsBuilder[T any] struct {
	Opts []func(*T) error
}

func (o *OptionsBuilder[T]) List() []func(*T) error {
	return o.Opts
}

func Option[T any]() *T {
	return new(T)
}
