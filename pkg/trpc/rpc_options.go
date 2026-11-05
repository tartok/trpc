package trpc

type OptionsBuilder[T any] struct {
	Opts []func(*T) error
}

func (o *OptionsBuilder[T]) List() []func(*T) error {
	return o.Opts
}

func Option[T any]() *T {
	return new(T)
}
