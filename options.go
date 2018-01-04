package gogo

import (
	"micro/codec"
	"micro/transport"
	"micro/registry"
	"micro/router"

	"context"
)

// Options of a service
type Options struct {
	Codec     codec.Codec
	Transport transport.Transport
	Registry  registry.Registry
	Router    router.Router

	//wrappers
	HdlrWrappers []HandlerWrapper
	//HttpHdlrWrappers []HttpHandlerWrapper

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Codec: codec.NewCodec(),
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Registry == nil {
		opt.Registry = registry.DefaultRegistry
	}

	if opt.Transport == nil {
		opt.Transport = transport.DefaultTransport
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

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// WrapHandler adds a service handler Wrapper to a list of options passed into the server
func WrapHandler(w ...HandlerWrapper) Option {
	return func(o *Options) {
		for _, wrap := range w {
			o.HdlrWrappers = append(o.HdlrWrappers, wrap)
		}
	}
}

// WrapHttpHandler adds a http handler Wrapper to a list of options passed into the server
//func WrapHttpHandler(w ...HttpHandlerWrapper) Option {
//	return func(o *Options) {
//		for _, wrap := range w {
//			o.HttpHdlrWrappers = append(o.HttpHdlrWrappers, wrap)
//		}
//	}
//}
