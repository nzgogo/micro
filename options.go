package gogo

import (
	"micro/transport"
)

// Options of a service
type Options struct {
	Transport transport.Transport
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	var opt Options

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		o.Transport.Init()
	}
}
