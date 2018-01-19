package gogo

import (
	"github.com/nzgogo/micro/context"
	"github.com/nzgogo/micro/registry"
	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/transport"
	"github.com/nzgogo/micro/selector"
)

// Options of a service
type Options struct {
	Transport transport.Transport
	Registry  registry.Registry
	Router    router.Router
	Context   context.Context
	Selector  selector.Selector

	//wrappers
	//HdlrWrappers []HandlerWrapper
	//HttpHdlrWrappers []HttpHandlerWrapper
	wrappers  wrapper
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.NewContext(),
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

func Context(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
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
