package trpc

type Lister[T any] interface {
	List() []func(*T) error
}

func AddOptions[T any](o *T, opts ...Lister[T]) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		for _, setArgs := range opt.List() {
			if setArgs == nil {
				continue
			}

			if err := setArgs(o); err != nil {
				return err
			}
		}
	}
	return nil
}
