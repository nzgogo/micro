package gogo

import (
	"net/http"

	"github.com/nzgogo/micro/router"
	"github.com/nzgogo/micro/codec"
)

// HttpHandlerFunc represents a single method of a http handler. It's used primarily
// for the wrappers.
type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
// HttpHandlerWrapper wraps the HttpHandlerFunc and returns the equivalent
type HttpHandlerWrapper func(HttpHandlerFunc) HttpHandlerFunc

// HandlerWrapper wraps the HandlerFunc and returns the equivalent
type HandlerWrapper func(router.Handler) router.Handler

//
type HttpResponseWriter func(rw http.ResponseWriter, response *codec.Message)
//
type HttpResponseWrapper func(HttpResponseWriter) HttpResponseWriter

// HttpWrapperChain builds the wrapper chain recursively, functions are first class
func HttpWrapperChain(f HttpHandlerFunc, m ...HttpHandlerWrapper) HttpHandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](HttpWrapperChain(f, m[1:cap(m)]...))
}
