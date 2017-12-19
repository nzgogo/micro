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
	opt := Options{
		Transport: transport.NewTransport(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
