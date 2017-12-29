package gogo
import (
	"micro/codec"
	"micro/transport"
)

// Options of a service
type Options struct {
	Codec     codec.Codec
	Transport transport.Transport

	HdlrWrappers []HandlerWrapper
	HttpHdlrWrappers []HttpHandlerWrapper
}

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

// WrapHandler adds a service handler Wrapper to a list of options passed into the server
func WrapHandler(w ...HandlerWrapper) Option {
	return func(o *Options) {
		for _, wrap := range w {
			o.HdlrWrappers = append(o.HdlrWrappers, wrap)
		}
	}
}

// WrapHttpHandler adds a http handler Wrapper to a list of options passed into the server
func WrapHttpHandler(w ...HttpHandlerWrapper) Option {
	return func(o *Options) {
		for _, wrap := range w {
			o.HttpHdlrWrappers = append(o.HttpHdlrWrappers, wrap)
		}
	}
}