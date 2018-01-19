package gogo

import (
	"net/http"

	"github.com/nzgogo/micro/codec"
	)

type wrapper struct{
	handlerWrappers []HandlerWrapper
	httpHandlerWrappers []HttpHandlerWrapper
}

// HttpHandlerFunc represents a single method of a http handler. It's used primarily
// for the wrappers.
type HttpHandlerFunc func(http.ResponseWriter, *http.Request) error

// HttpHandlerWrapper wraps the HttpHandlerFunc and returns the equivalent
type HttpHandlerWrapper func(HttpHandlerFunc) HttpHandlerFunc

// HandlerFunc represents a single method of a service router handler. It's used primarily
// for the wrappers (after api interpreter and before service handler).
type HandlerFunc func(*codec.Message) error

// HandlerWrapper wraps the HandlerFunc and returns the equivalent
type HandlerWrapper func(HandlerFunc) HandlerFunc

