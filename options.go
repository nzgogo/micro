package gogo
import (
	"micro/codec"
	"micro/transport"
)

// Options of a service
type Options struct {
	Codec     codec.Codec
	Transport transport.Transport
}

type Option func(*Options)

func newOptions(opts ...Option) Options {

	opt := Options{
		Codec: codec.NewCodec(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}
